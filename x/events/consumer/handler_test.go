package consumer

import (
	"context"
	"fmt"
)

func newHandler(name string) Handler {
	return func(ctx context.Context, m Message) error {
		fmt.Println(name)
		return nil
	}
}

func newMiddleware(name string) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, m Message) error {
			fmt.Println(name)
			return next(ctx, m)
		}
	}
}

func ExampleApplyMiddlewares() {

	handler := ApplyMiddlewares(
		newHandler("handler"),
		[]Middleware{
			newMiddleware("mdw-1"),
			newMiddleware("mdw-2"),
			newMiddleware("mdw-3"),
		},
	)

	handler(nil, nil)

	// Output:
	// mdw-1
	// mdw-2
	// mdw-3
	// handler
}
