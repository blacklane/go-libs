package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/blacklane/go-libs/graceful"
)

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello World!"))
	})

	srv := &http.Server{
		Handler: m,
		Addr:    ":8081",
	}

	g := graceful.New(
		graceful.WithBeforeStartHooks(func(ctx context.Context) error {
			log.Println("open db connection")
			return nil
		}),
		graceful.WithAfterStopHooks(func(ctx context.Context) error {
			log.Println("closing db connection")
			return nil
		}),
		graceful.WithTasks(
			// http server
			graceful.NewHTTPServerTask(srv),
			// interval task
			graceful.NewIntervalTask(2*time.Second, func(ctx context.Context) error {
				log.Println("interval task!")
				return nil
			}),
			// custom start/stop functions
			graceful.NewTask(
				func(ctx context.Context) error {
					log.Println("custom server - start")
					return nil
				},
				func(ctx context.Context) error {
					log.Println("custom server - stop")
					return nil
				},
			),
		),
	)

	// appending task
	g.AppendTask(graceful.NewTask(
		func(ctx context.Context) error {
			log.Println("appended custom server - start")
			return nil
		},
		func(ctx context.Context) error {
			log.Println("appended custom server - stop")
			return nil
		},
	))

	log.Printf("server running %s\n", srv.Addr)

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
