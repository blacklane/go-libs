module github.com/blacklane/go-libs/otel

go 1.16

require (
	github.com/blacklane/go-libs/logger v0.2.0
	github.com/blacklane/go-libs/tracking v0.2.1
	github.com/blacklane/go-libs/x/events v0.0.5
	github.com/google/uuid v1.3.0
	github.com/rs/zerolog v1.24.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC3
	go.opentelemetry.io/otel/sdk v1.0.0-RC3
	go.opentelemetry.io/otel/trace v1.0.0-RC3
)

replace github.com/blacklane/go-libs/x/events => ../x/events

replace github.com/blacklane/go-libs/logger => ../logger
