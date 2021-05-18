module github.com/blacklane/go-libs/tracking/examples

go 1.16

require (
	github.com/blacklane/go-libs/logger v0.5.2-0.20210512094936-44b87e762488
	github.com/blacklane/go-libs/tracking v0.2.1
	github.com/blacklane/go-libs/x/events v0.1.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/google/uuid v1.2.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/rs/zerolog v1.22.0 // indirect
)

replace github.com/blacklane/go-libs/tracking => ../
