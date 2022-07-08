package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func TestFromContext(t *testing.T) {
	const (
		message = "Hello!"
		appName = "test"
	)
	var buf bytes.Buffer
	logger := New(&buf, appName)
	ctx := logger.WithContext(context.Background())
	log := FromContext(ctx)

	log.Info().Msg(message)
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	if data["message"] != message {
		t.Errorf("expected message to be '%s', got %q", message, data["message"])
	}
	if data["application"] != appName {
		t.Errorf("expected application to be '%s', got %q", appName, data["application"])
	}
}

func TestFromContext_TraceContext(t *testing.T) {
	traceID := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}

	ddTraceID := "651345242494996240"
	ddSpanID := "72623859790382856"

	var buf bytes.Buffer
	logger := New(&buf, "test")
	ctx := logger.WithContext(context.Background())

	span := newFakeSpan(traceID, spanID)

	ctx = trace.ContextWithSpan(ctx, span)

	log := FromContext(ctx)

	log.Info().Msg("message")
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	if data["dd.trace_id"] != ddTraceID {
		t.Errorf("expected dd.trace_id to be '%s', got %q", ddTraceID, data["dd.trace_id"])
	}
	if data["dd.span_id"] != ddSpanID {
		t.Errorf("expected dd.span_id to be '%s', got %q", ddSpanID, data["dd.span_id"])
	}
}

type fakeSpan struct {
	Context trace.SpanContext
}

func newFakeSpan(traceID [16]byte, spanID [8]byte) trace.Span {
	ctx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})

	return fakeSpan{
		Context: ctx,
	}
}

func (f fakeSpan) SpanContext() trace.SpanContext        { return f.Context }
func (fakeSpan) IsRecording() bool                       { return true }
func (fakeSpan) SetStatus(codes.Code, string)            {}
func (fakeSpan) SetError(bool)                           {}
func (fakeSpan) SetAttributes(...attribute.KeyValue)     {}
func (fakeSpan) End(...trace.SpanEndOption)              {}
func (fakeSpan) RecordError(error, ...trace.EventOption) {}
func (fakeSpan) AddEvent(string, ...trace.EventOption)   {}
func (fakeSpan) SetName(string)                          {}
func (fakeSpan) TracerProvider() trace.TracerProvider    { return nil }
