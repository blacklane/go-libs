package probes

import "context"

type ctxK int

const ctxKey = ctxK(42)

// FromCtx returns the Probe associated with the context.
// A noop Probe is returned if there is no Probe associated with the context.
func FromCtx(ctx context.Context) Probe {
	val, ok := ctx.Value(ctxKey).(Probe)
	if !ok {
		return NewNoop()
	}

	return val
}
