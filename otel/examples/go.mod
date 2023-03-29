module github.com/blacklane/go-libs/otel/examples

go 1.17

replace (
	github.com/blacklane/go-libs/camunda/v2 => ../../camunda/v2
	github.com/blacklane/go-libs/logger => ../../logger
	github.com/blacklane/go-libs/middleware => ../../middleware
	github.com/blacklane/go-libs/otel => ../
	github.com/blacklane/go-libs/tracking => ../../tracking
	github.com/blacklane/go-libs/x/events => ../../x/events
)

require (
	github.com/blacklane/go-libs/logger v0.7.2
	github.com/blacklane/go-libs/middleware v0.1.0
	github.com/blacklane/go-libs/otel v0.2.3
	github.com/blacklane/go-libs/tracking v0.3.1
	github.com/blacklane/go-libs/x/events v0.5.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/google/uuid v1.3.0
	github.com/rs/zerolog v1.29.0
	go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/otel/trace v1.14.0
)

require (
	github.com/blacklane/go-libs/camunda/v2 v2.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/proto/otlp v0.16.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220503193339-ba3ae3f07e29 // indirect
	google.golang.org/grpc v1.46.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
