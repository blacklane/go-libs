package events

import "context"

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

func (h HandlerFunc) Handle(ctx context.Context, e Event) error {
	return h(ctx, e)
}

type HandlerBuilder struct {
	middleware  []Middleware
	rawHandlers []Handler
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
		for _, m := range hb.middleware {
			h = m(h)
		}
		handlers = append(handlers, h)
	}
	return handlers
}
