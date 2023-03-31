package hooks

import (
	"context"
	"strconv"

	"github.com/blacklane/go-libs/logger/internal"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

const (
	ddTraceIDKey = "dd.trace_id"
	ddSpanIDKey  = "dd.span_id"
)

func TraceContext(ctx context.Context) zerolog.Hook {
	return &traceContextHook{ctx: ctx}
}

type traceContextHook struct {
	ctx context.Context
}

func (h *traceContextHook) Run(ev *zerolog.Event, level zerolog.Level, msg string) {
	span := trace.SpanFromContext(h.ctx)
	if span.IsRecording() {
		spanCtx := span.SpanContext()
		traceID := spanCtx.TraceID().String()
		spanID := spanCtx.SpanID().String()

		ddTraceID := convertTraceID(traceID)
		ddSpanID := convertTraceID(spanID)

		ev.
			Str(ddTraceIDKey, ddTraceID).
			Str(ddSpanIDKey, ddSpanID).
			Str(internal.FieldTraceID, ddTraceID)
	}
}

// convertTraceID translates OpenTelemetry formatted trace_id and span_id into the Datadog format.
// https://docs.datadoghq.com/tracing/other_telemetry/connect_logs_and_traces/opentelemetry/?tab=go
func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
