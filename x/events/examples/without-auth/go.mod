module github.com/blacklane/go-libs/x/events/examples/without-auth

go 1.17

replace (
	github.com/blacklane/go-libs/tracking => ../../../../tracking
	github.com/blacklane/go-libs/x/events => ../../
)

require (
	github.com/blacklane/go-libs/x/events v0.0.0-20200729083625-bff51d3ec664
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/confluentinc/confluent-kafka-go v1.9.1
)

require (
	github.com/blacklane/go-libs/tracking v0.3.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
