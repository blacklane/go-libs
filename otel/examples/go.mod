module github.com/blacklane/go-libs/otel/examples

go 1.17

require (
	github.com/blacklane/go-libs/logger v0.6.0
	github.com/blacklane/go-libs/middleware v0.1.0
	github.com/blacklane/go-libs/otel v0.1.0
	github.com/blacklane/go-libs/tracking v0.3.0
	github.com/blacklane/go-libs/x/events v0.2.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/google/uuid v1.3.0
	github.com/rs/zerolog v1.26.0
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/trace v1.1.0
)

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.26.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.1.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.1.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.1.0 // indirect
	go.opentelemetry.io/proto/otlp v0.10.0 // indirect
	golang.org/x/net v0.0.0-20211109214657-ef0fda0de508 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.0.0-20211109184856-51b60fd695b3 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211104193956-4c6863e31247 // indirect
	google.golang.org/grpc v1.42.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/blacklane/go-libs/otel => ../
