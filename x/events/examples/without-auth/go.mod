module github.com/blacklane/go-libs/x/events/examples/without-auth

go 1.17

replace github.com/blacklane/go-libs/x/events => ../../

require (
	github.com/blacklane/go-libs/x/events v0.0.0-20200729083625-bff51d3ec664
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/confluentinc/confluent-kafka-go v1.7.0
)

require (
	github.com/blacklane/go-libs/tracking v0.2.1 // indirect
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	go.opentelemetry.io/otel v1.0.0-RC2 // indirect
	go.opentelemetry.io/otel/trace v1.0.0-RC2 // indirect
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	google.golang.org/appengine v1.4.0 // indirect
)
