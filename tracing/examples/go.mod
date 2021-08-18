module github.com/blacklane/go-libs/tracing/examples

go 1.16

require (
	github.com/blacklane/go-libs/logger v0.5.1
	github.com/blacklane/go-libs/tracing v0.0.0
	github.com/blacklane/go-libs/tracking v0.2.1
	github.com/blacklane/go-libs/x/events v0.1.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/google/uuid v1.2.0
	github.com/rs/zerolog v1.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/trace v1.0.0-RC2
)

replace github.com/blacklane/go-libs/tracing => ../

replace github.com/blacklane/go-libs/x/events => ../../x/events

replace github.com/blacklane/go-libs/logger => ../../logger
