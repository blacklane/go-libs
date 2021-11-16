package events

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

type Header map[string]string
type Event struct {
	Headers Header
	Key     []byte
	Payload []byte
}

type Middleware func(Handler) Handler

type Handler interface {
	Handle(context.Context, Event) error
}

type HandlerFunc func(ctx context.Context, e Event) error

type HandlerBuilder struct {
	middleware  []Middleware
	rawHandlers []Handler
}

// Ensure Handler implements OTel propagation.TextMapCarrier
var _ = propagation.TextMapCarrier(Header{})

func (h Header) Get(key string) string {
	return h[key]
}

func (h Header) Set(key, value string) {
	h[key] = value
}

func (h Header) Keys() []string {
	keys := make([]string, len(h))
	for k := range h {
		keys = append(keys, k)
	}

	return keys
}

func (h HandlerFunc) Handle(ctx context.Context, e Event) error {
	return h(ctx, e)
}

func (hb *HandlerBuilder) UseMiddleware(m ...Middleware) {
	hb.middleware = append(hb.middleware, m...)
}

func (hb *HandlerBuilder) AddHandler(h Handler) {
	hb.rawHandlers = append(hb.rawHandlers, h)
}

func (hb HandlerBuilder) Build() []Handler {
	var handlers []Handler
	for _, rh := range hb.rawHandlers {
		h := rh
		for i := len(hb.middleware) - 1; i >= 0; i-- {
			h = hb.middleware[i](h)
		}
		handlers = append(handlers, h)
	}
	return handlers
}
