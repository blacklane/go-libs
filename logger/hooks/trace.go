package hooks

import (
	"context"
	"strconv"

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

func (ot *traceContextHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	span := trace.SpanFromContext(ot.ctx)
	if span.IsRecording() {
		spanCtx := span.SpanContext()
		traceID := spanCtx.TraceID().String()
		spanID := spanCtx.SpanID().String()

		e.Str(ddTraceIDKey, convertTraceID(traceID)).
			Str(ddSpanIDKey, convertTraceID(spanID))
	}
}

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
