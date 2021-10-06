package events

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/blacklane/go-libs/tracking"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var ErrProducerNotHandlingEvents = errors.New("producer should be handling events")
var ErrProducerIsAlreadyRunning = errors.New("producer is already running")

type Producer interface {
	// Send an event to the given topic
	// deprecated, use SendCtx instead
	Send(event Event, topic string) error
	// SendCtx
	// TODO(Add docs):
	SendCtx(ctx context.Context, eventName string, event Event, topic string) error
	// SendWithTrackingID adds the tracking ID to the event's headers and sends it to the given topic
	SendWithTrackingID(trackingID string, event Event, topic string) error
	// HandleEvents starts to listen to the producer events channel
	HandleEvents() error
	// Shutdown gracefully shuts down the producer, it respect the context
	// timeout.
	Shutdown(ctx context.Context) error
}

// KafkaProducerConfig holds all possible configurations for the kafka producer.
// Use NewKafkaProducerConfig to initialise it.
// To see the possible configurations, check the its WithXXX methods and
// *kafkaConfig.WithXXX` methods as well
type KafkaProducerConfig struct {
	*kafkaConfig

	deliveryErrHandler func(Event, error)

	flushTimeoutMs int
}

type kafkaProducer struct {
	*kafkaCommon

	config   *kafka.ConfigMap
	producer *kafka.Producer

	deliveryErrHandler func(Event, error)

	flushTimeoutMs int

	otelPropagator propagation.TextMapPropagator

	runMu    sync.RWMutex
	run      bool
	shutdown bool
}

// NewKafkaProducerConfig returns an initialised *KafkaProducerConfig
func NewKafkaProducerConfig(config *kafka.ConfigMap) *KafkaProducerConfig {
	return &KafkaProducerConfig{
		kafkaConfig: &kafkaConfig{
			config:      config,
			tokenSource: emptyTokenSource{},
			errFn:       func(error) {},
		},
		deliveryErrHandler: func(Event, error) {},
	}
}

// WithEventDeliveryErrHandler registers a delivery error handler to be called
// whenever a delivery fails.
func (pc *KafkaProducerConfig) WithEventDeliveryErrHandler(errHandler func(Event, error)) {
	pc.deliveryErrHandler = errHandler
}

// WithFlushTimeout sets the producer Flush timeout.
func (pc *KafkaProducerConfig) WithFlushTimeout(timeout int) {
	pc.flushTimeoutMs = timeout
}

// NewKafkaProducer returns new a producer.
// To handle errors, either `kafka.Error` messages or any other error while
// interacting with Kafka, register an Error function on *KafkaConsumerConfig.
func NewKafkaProducer(c *KafkaProducerConfig) (Producer, error) {
	kp := &kafkaProducer{
		kafkaCommon: &kafkaCommon{
			errFn:       c.errFn,
			tokenSource: c.tokenSource,
		},
		config: c.config,

		deliveryErrHandler: c.deliveryErrHandler,
		flushTimeoutMs:     c.flushTimeoutMs,

		otelPropagator: otel.GetTextMapPropagator(),
	}

	p, err := kafka.NewProducer(c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafkaProducer: %w", err)
	}
	kp.producer = p

	return kp, nil
}

// HandleEvents listens to the producer events channel.
// To register handlers check *KafkaProducerConfig.
func (p *kafkaProducer) HandleEvents() error {
	if p.running() {
		return ErrProducerIsAlreadyRunning
	}
	p.startRunning()

	go func() {
		for kafkaEvent := range p.producer.Events() {
			switch kafkaEvent := kafkaEvent.(type) {
			case *kafka.Message:
				if kafkaEvent.TopicPartition.Error != nil {
					if p.deliveryErrHandler != nil {
						event := messageToEvent(kafkaEvent)
						p.deliveryErrHandler(
							*event,
							fmt.Errorf(
								"%w, topic '%s', partition: '%d', offset: '%s'",
								kafkaEvent.TopicPartition.Error,
								*kafkaEvent.TopicPartition.Topic,
								kafkaEvent.TopicPartition.Partition,
								kafkaEvent.TopicPartition.Offset.String()))
					}
				}
			case kafka.OAuthBearerTokenRefresh:
				p.refreshToken()
			case kafka.Error:
				p.errFn(kafkaEvent)
			}
		}
	}()

	return nil
}

func (p *kafkaProducer) refreshToken() {
	p.kafkaCommon.refreshToken(p.producer)
}

// SendCtx adds the tracking ID to the event's headers. If the event already
// has the tracking ID header set, it does nothing.
func (p *kafkaProducer) SendCtx(ctx context.Context, eventName string, event Event, topic string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("will not send event: %w", err)
	}

	ctx, sp := otel.Tracer(OTelTracerName).
		Start(ctx,
			"produced:"+eventName,
			trace.WithSpanKind(trace.SpanKindProducer),
			trace.WithAttributes(semconv.MessagingDestinationKey.String(topic)))
	defer sp.End()

	sp.SetAttributes(semconv.MessagingMessageIDKey.String(string(event.Key)))

	if event.Headers == nil {
		event.Headers = map[string]string{}
	}
	p.otelPropagator.Inject(ctx, event.Headers)

	trackedEvent := addTrackingID(tracking.IDFromContext(ctx), event)

	return p.Send(trackedEvent, topic)
}

// SendWithTrackingID adds the tracking ID to the event's headers. If the event already
// has the tracking ID header set, it does nothing.
func (p *kafkaProducer) SendWithTrackingID(trackingID string, event Event, topic string) error {
	trackedEvent := addTrackingID(trackingID, event)
	return p.Send(trackedEvent, topic)
}

// Send messages to the given topic. Delivery errors are sent to the producer's
// event channel. To handle delivery errors check *KafkaProducerConfig.WithEventDeliveryErrHandler.
func (p *kafkaProducer) Send(event Event, topic string) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            event.Key,
		Value:          event.Payload,
		Headers:        toKafkaHeaders(event.Headers),
	}, nil)
}

// Shutdown the producer and also closes the underlying kafka producer.
func (p *kafkaProducer) Shutdown(ctx context.Context) error {
	isRunning := p.running()
	p.stopRun()
	defer p.producer.Close()

	// Flush will only work if we listen to producer's events
	if !isRunning {
		return ErrProducerNotHandlingEvents
	}

	for {
		notSent := p.producer.Flush(p.flushTimeoutMs)
		if notSent == 0 {
			break
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%d messages not sent", notSent)
		default:
			continue
		}
	}

	return nil
}

func (p *kafkaProducer) startRunning() bool {
	p.runMu.RLock()
	p.run = true
	defer p.runMu.RUnlock()
	return p.run
}

func (p *kafkaProducer) running() bool {
	p.runMu.RLock()
	defer p.runMu.RUnlock()
	return p.run
}

func (p *kafkaProducer) stopRun() {
	p.runMu.Lock()
	defer p.runMu.Unlock()
	p.run = false
}

func toKafkaHeaders(eventHeaders Header) []kafka.Header {
	var kafkaHeaders []kafka.Header
	for key, value := range eventHeaders {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: key, Value: []byte(value)})
	}
	return kafkaHeaders
}

func addTrackingID(trackingID string, event Event) Event {
	if trackingID != "" {
		if event.Headers == nil || len(event.Headers) == 0 {
			event.Headers = Header{}
		}

		if event.Headers[HeaderTrackingID] == "" {
			event.Headers[HeaderTrackingID] = trackingID
		}
	}

	return event
}
