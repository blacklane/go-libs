package events

import (
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestKafkaProducer_WithOAuth(t *testing.T) {
	kpc := NewKafkaProducerConfig(&kafka.ConfigMap{})
	tokenSource := tokenSourceMock{}
	kpc.WithOAuth(tokenSource)

	p, err := NewKafkaProducer(kpc)
	if err != nil {
		t.Fatalf("could not create kafka consumer: %v", err)
	}

	kp, ok := p.(*kafkaProducer)
	if !ok {
		t.Fatalf("cannot cast Consumer to *kafkaProducer")
	}

	if !cmp.Equal(kp.tokenSource, tokenSource) {
		t.Errorf("want: %v, got: %v", tokenSource, kp.tokenSource)
	}
}

func TestKafkaProducer_WithErrFunc(t *testing.T) {
	var got error
	want := errors.New("error TestKafkaConfig_WithErrFunc")
	errFn := func(err error) { got = err }

	kpc := NewKafkaProducerConfig(&kafka.ConfigMap{"group.id": "TestKafkaConsumer_WithErrFunc"})
	kpc.WithErrFunc(errFn)

	p, err := NewKafkaProducer(kpc)
	if err != nil {
		t.Fatalf("cannot create kafka consumer: %v", err)
	}

	kp, ok := p.(*kafkaProducer)
	if !ok {
		t.Fatalf("cannot cast Consumer to *kafkaProducer")
	}

	kp.errFn(want)
	if got != want {
		t.Errorf("want: %v, got %v", want, got)
	}
}

func TestNewKafkaProducerConfigAllInitialised(t *testing.T) {
	kc := NewKafkaProducerConfig(nil)

	if kc.deliveryErrHandler == nil {
		t.Errorf("deliveryErrHandler is nil")
	}
	if kc.tokenSource == nil {
		t.Errorf("tokenSource is nil")
	}
	if kc.errFn == nil {
		t.Errorf("errFn is nil")
	}
}

func TestParseProducerHeaders(t *testing.T) {
	key, value := "key", "value"
	eventHeaders := Header{
		key: value,
	}

	kafkaHeaders := parseProducerHeaders(eventHeaders)

	assert.Equal(t, eventHeaders, parseHeaders(kafkaHeaders))
}
