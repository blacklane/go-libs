package consumer

import "context"

type Handler func(ctx context.Context, m Message) error

type Middleware func(next Handler) Handler

func ApplyMiddlewares(handler Handler, middlewares []Middleware) Handler {
	h := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
