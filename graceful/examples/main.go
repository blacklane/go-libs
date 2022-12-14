package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/blacklane/go-libs/graceful"
)

func IntervalTask(ctx context.Context) error {
	log.Println("interval task!")
	return nil
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	srv := &http.Server{
		Handler: m,
		Addr:    ":8081",
	}

	g := graceful.New(
		graceful.WithServers(
			// http server
			graceful.NewHTTPServer(srv),
			// interval task
			graceful.NewIntervalTaskServer(2*time.Second, IntervalTask),
			// custom start/stop functions
			graceful.NewServer(
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

	log.Printf("server running %s\n", srv.Addr)

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
