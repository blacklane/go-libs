package examples

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/uhttp"
	"github.com/blacklane/go-libs/x/events"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/blacklane/go-libs/otel"
)

// StartHTTPServer creates and starts a http.Server listening on port 4242, with no router
// and a single handler. See newHandler for details about the handler.
func StartHTTPServer(serviceName string, producer events.Producer, topic string, eventName string) {
	log := logger.New(logger.ConsoleWriter{Out: os.Stdout}, serviceName).With().
		Str("environment", "otel").Logger()

	path := "/tracing/example/path"

	// Using HTTPAllMiddleware as it's been applied directly to the handler.
	applyMiddleware := uhttp.HTTPAllMiddleware(
		serviceName,
		"my-awesome-handler",
		path,
		log)

	handler := applyMiddleware(newHandler(producer, topic, eventName))

	httpServer := http.Server{
		Addr:    ":4242",
		Handler: handler,
	}

	go func() { httpServer.ListenAndServe() }()
	log.Info().Msgf("Starting HTTP server on %s", httpServer.Addr)
}

func newHandler(producer events.Producer, topic string, eventName string) http.Handler {
	var count int
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sp := trace.SpanFromContext(ctx)
		sp.SetAttributes(
			attribute.String("some_uuid", uuid.New().String()),
			attribute.String("some_key", "some_value"),
			attribute.Int("some_int", count))

		defer func() { count++ }()

		// The headers will be sent as part of the response body to show the
		// headers OpenTelemetry uses.
		headers, _ := json.Marshal(w.Header())

		// simulates a failure by flipping a coin.
		if rand.Int()%2 == 0 {
			err := errors.New(http.StatusText(http.StatusTeapot))

			otel.SpanAddErr(sp, err)
			logger.FromContext(ctx).Err(err).Msg("handler failed: bad luck")

			w.WriteHeader(http.StatusTeapot)
			_, _ = fmt.Fprintf(w, "I'm a tea pot\ntracking_id: %s\nheaders: %s",
				tracking.IDFromContext(ctx),
				headers)
			return
		}

		time.Sleep(time.Duration(5*rand.Intn(5)) * time.Millisecond)
		err := produceEvent(ctx, producer, topic, eventName, count)
		if err != nil {
			err := fmt.Errorf("could not send event: %w", err)

			otel.SpanAddErr(sp, err)
			logger.FromContext(ctx).Err(err).
				Str("event", eventName).
				Msg("handler failed to produce event")

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

func produceEvent(ctx context.Context, producer events.Producer, topic string, eventName string, count int) error {
	event := events.Event{
		Headers: map[string]string{"X-Example": "a-value"},
		Key:     []byte(fmt.Sprintf("%d", count)),
		Payload: []byte(fmt.Sprintf(`{"event":"%s","count":%d}`, eventName, count)),
	}

	logger.FromContext(ctx).Debug().
		Str("event_produced_headers", fmt.Sprintf("%v", event.Headers)).
		RawJSON("event_produced_payload", event.Payload).
		Str("event_produced_topic", topic).Msg("producing event")

	return producer.SendCtx(ctx, eventName, event, topic)
}
