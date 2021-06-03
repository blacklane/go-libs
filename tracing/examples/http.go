package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blacklane/go-libs/tracing"
)

func newHTTPServerWithOTel(serviceName string, producer events.Producer) {
	// Creates a logger for this "service"
	log := logger.New(logger.ConsoleWriter{Out: os.Stdout}, serviceName)

	path := "/tracing/example/path"

	middleware := tracing.HTTPDefaultOTelMiddleware(serviceName, path, log)
	handler := middleware(httpHandler(producer, serviceName))

	httpServer := http.Server{
		Addr:    ":4242",
		Handler: handler,
	}

	go func() { httpServer.ListenAndServe() }()
	log.Info().Msgf("Starting HTTP server on %s", httpServer.Addr)
}

func httpHandler(producer events.Producer, serviceName string) http.Handler {
	var count int
	propagator := otel.GetTextMapPropagator()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glTracer := otel.GetTracerProvider()
		tracer := glTracer.Tracer(serviceName, trace.WithInstrumentationVersion("development"))
		ctx, sp := tracer.
			Start(r.Context(),
				"http_handler", // Span name
				trace.WithAttributes(
					attribute.String("some_uuid", uuid.New().String()),
					attribute.String("some_key", "some_value"),
					attribute.Int("some_int", count)))
		defer sp.End()

		defer func() { count++ }()

		// The headers will be sent as part of the response body to show the
		// headers OpenTelemetry uses.
		headers, _ := json.Marshal(w.Header())

		// simulates a failure by flipping a coin.
		if rand.Int()%2 == 0 {
			err := errors.New(http.StatusText(http.StatusTeapot))
			sp.RecordError(err)

			w.WriteHeader(http.StatusTeapot)
			_, _ = fmt.Fprintf(w, "I'm a tea pot\ntracking_id: %s\nheaders: %s",
				tracking.IDFromContext(ctx),
				headers)
			return
		}

		event := events.Event{
			Headers: map[string]string{},
			Key:     []byte(fmt.Sprintf("%d", count)),
			Payload: []byte(fmt.Sprintf(`{"event":"%s","count":%d}`, eventName, count)),
		}

		// logger.FromContext(ctx).Debug().
		// 	Str("ctx", fmt.Sprintf("%v", ctx)).
		// 	Msg("ctx before propagator.Inject")
		// logger.FromContext(ctx).Debug().
		// 	Str("event_headers", fmt.Sprintf("%v", event.Headers)).
		// 	Msg("before propagator.Inject")
		propagator.Inject(ctx, event.Headers)
		// logger.FromContext(ctx).Debug().
		// 	Str("event_headers", fmt.Sprintf("%v", event.Headers)).
		// 	Msg("after propagator.Inject")

		logger.FromContext(ctx).Debug().
			Str("event_produced_headers", fmt.Sprintf("%v", event.Headers)).
			RawJSON("event_produced_payload", event.Payload).
			Str("event_produced_topic", topic).Msg("producing event")
		err := producer.Send(event, topic)
		if err != nil {
			err := fmt.Errorf("could not send event: %w", err)

			sp.RecordError(err)

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"error":"internal server error","tracking_id":"%s","headers":%q}`,
				tracking.IDFromContext(ctx),
				headers)
			return
		}

		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(
			fmt.Sprintf(`{"tracking_id":"%s","headers":%q}`,
				tracking.IDFromContext(ctx),
				headers)))
	})
}
