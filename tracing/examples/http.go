package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/blacklane/go-libs/x/events"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracinglog "github.com/opentracing/opentracing-go/log"

	"github.com/blacklane/go-libs/tracing"
	"github.com/blacklane/go-libs/tracing/internal/constants"
)

func newHTTPServer(serviceName string, tracer opentracing.Tracer, producer events.Producer) {
	// Creates a logger for this "service"
	log := logger.New(logger.ConsoleWriter{Out: os.Stdout}, serviceName)

	path := "/tracing/example/path"

	middleware := tracing.HTTPDefaultMiddleware(path, tracer, log)
	handler := middleware(httpHandler(producer))
	httpServer := http.Server{
		Addr:    ":4242",
		Handler: handler,
	}

	go func() { httpServer.ListenAndServe() }()
	log.Info().Msgf("Starting HTTP server on %s", httpServer.Addr)
}

func httpHandler(producer events.Producer) http.Handler {
	var count int
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { count++ }()
		ctx := r.Context()

		dump, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(dump))

		// Gets the opentracing span
		sp := tracing.SpanFromContext(ctx)

		// add fields to a log entry on this span
		sp.LogFields(
			opentracinglog.String("some_uuid", uuid.New().String()),
			opentracinglog.String("some_key", "some_value"),
			opentracinglog.Int("some_number", 1500),
			opentracinglog.Message("some message"))

		headers, _ := json.Marshal(w.Header())

		// simulates a failure by flipping a coin.
		if rand.Int()%2 == 0 {
			err := errors.New(http.StatusText(http.StatusTeapot))

			// Flags an error happened on this span
			ext.Error.Set(sp, true)

			// Logs the error on the span
			sp.LogFields(opentracinglog.Error(err))

			w.WriteHeader(http.StatusTeapot)
			_, _ = fmt.Fprintf(w, "I'm a tea pot\ntracking_id: %s\nheaders: %s",
				tracking.IDFromContext(ctx),
				headers)
			return
		}

		event := events.Event{
			Headers: nil,
			Key:     []byte(fmt.Sprintf("%d", count)),
			Payload: []byte(fmt.Sprintf(`{"event":"%s","count":%d}`, eventName, count)),
		}

		// Injects the span into the events headers.
		// It's how the span is propagated through different services.
		err := tracing.EventsOpentracingInject(ctx, sp, &event)
		if err != nil {
			logger.FromContext(ctx).Err(err).
				Str(constants.FieldEvent, eventName).
				Msg("could not inject span into event")
		}

		err = producer.Send(event, topic)
		if err != nil {
			err := fmt.Errorf("could not send event: %w", err)

			// Flags an error happened on this span
			ext.Error.Set(sp, true)

			// Logs the error on the span
			sp.LogFields(opentracinglog.Error(err))

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
