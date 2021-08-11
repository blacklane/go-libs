module github.com/blacklane/go-libs/tracing

go 1.16

require (
	github.com/blacklane/go-libs/logger v0.5.1
	github.com/blacklane/go-libs/tracking v0.2.1
	github.com/blacklane/go-libs/x/events v0.1.0
	github.com/confluentinc/confluent-kafka-go v1.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.2.0
	github.com/rs/zerolog v1.22.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/otlp v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c // indirect
)

replace github.com/blacklane/go-libs/x/events => ../x/events

replace github.com/blacklane/go-libs/logger => ../logger
