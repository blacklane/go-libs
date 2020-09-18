package events

import (
	"context"
	"fmt"
)

func ExampleHandlerBuilder() {
	m1 := Middleware(func(handler Handler) Handler {
		return HandlerFunc(func(ctx context.Context, e Event) error {
			fmt.Println("middleware 1: before handler")
			err := handler.Handle(ctx, e)
			fmt.Println("middleware 1: after handler")
			return err
		})
	})
	m2 := Middleware(func(handler Handler) Handler {
		return HandlerFunc(func(ctx context.Context, e Event) error {
			fmt.Println("middleware 2: before handler")
			err := handler.Handle(ctx, e)
			fmt.Println("middleware 2: after handler")
			return err
		})
	})
	h := HandlerFunc(func(_ context.Context, _ Event) error {
		fmt.Println("handler")
		return nil
	})

	hb := HandlerBuilder{}
	hb.AddHandler(h)
	hb.UseMiddleware(m1, m2)

	// HandlerBuilder.Build returns a slice as several handlers might be added
	handler := hb.Build()[0]

	err := handler.Handle(context.Background(), Event{})
	if err != nil {
		fmt.Println("handler error: " + err.Error())
	}

	// Output:
	// middleware 1: before handler
	// middleware 2: before handler
	// handler
	// middleware 2: after handler
	// middleware 1: after handler
}

func ExampleHandlerBuilder_multipleHandlers() {
	m1 := Middleware(func(handler Handler) Handler {
		return HandlerFunc(func(ctx context.Context, e Event) error {
			fmt.Println("middleware 1: before handler")
			err := handler.Handle(ctx, e)
			fmt.Println("middleware 1: after handler")
			return err
		})
	})
	m2 := Middleware(func(handler Handler) Handler {
		return HandlerFunc(func(ctx context.Context, e Event) error {
			fmt.Println("middleware 2: before handler")
			err := handler.Handle(ctx, e)
			fmt.Println("middleware 2: after handler")
			return err
		})
	})
	h1 := HandlerFunc(func(_ context.Context, _ Event) error {
		fmt.Println("handler 1")
		return nil
	})
	h2 := HandlerFunc(func(_ context.Context, _ Event) error {
		fmt.Println("handler 2")
		return nil
	})

	hb := HandlerBuilder{}
	hb.AddHandler(h1)
	hb.UseMiddleware(m1, m2)
	hb.AddHandler(h2)

	// HandlerBuilder.Build returns a slice as several handlers might be added
	handlers := hb.Build()
	handler1 := handlers[0]
	handler2 := handlers[1]

	err := handler1.Handle(context.Background(), Event{})
	if err != nil {
		fmt.Println("handler1 error: " + err.Error())
	}

	fmt.Print("\n")

	err = handler2.Handle(context.Background(), Event{})
	if err != nil {
		fmt.Println("handler2 error: " + err.Error())
	}

	// Output:
	// middleware 1: before handler
	// middleware 2: before handler
	// handler 1
	// middleware 2: after handler
	// middleware 1: after handler
	//
	// middleware 1: before handler
	// middleware 2: before handler
	// handler 2
	// middleware 2: after handler
	// middleware 1: after handler
}
