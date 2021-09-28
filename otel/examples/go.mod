module github.com/blacklane/go-libs/otel/examples

go 1.17

require (
	github.com/blacklane/go-libs/logger v0.5.1
	github.com/blacklane/go-libs/middleware v0.0.0-00010101000000-000000000000
	github.com/blacklane/go-libs/otel v0.0.0
	github.com/blacklane/go-libs/tracking v0.2.1
	github.com/blacklane/go-libs/x/events v0.1.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/google/uuid v1.3.0
	github.com/rs/zerolog v1.24.0
	go.opentelemetry.io/otel v1.0.0
	go.opentelemetry.io/otel/trace v1.0.0
)

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.0.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.23.0 // indirect
	go.opentelemetry.io/otel/metric v0.23.0 // indirect
	go.opentelemetry.io/otel/sdk v1.0.0 // indirect
	go.opentelemetry.io/proto/otlp v0.9.0 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/blacklane/go-libs/otel => ../

replace github.com/blacklane/go-libs/x/events => ../../x/events

replace github.com/blacklane/go-libs/logger => ../../logger

replace github.com/blacklane/go-libs/middleware => ../../middleware
