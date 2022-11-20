module github.com/blacklane/go-libs/x/events

go 1.17

replace github.com/blacklane/go-libs/tracking => ../../tracking

require (
	github.com/blacklane/go-libs/tracking v0.3.1
	github.com/confluentinc/confluent-kafka-go v1.9.1
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.9
	github.com/google/uuid v1.3.0
	go.opentelemetry.io/otel v1.11.1
	go.opentelemetry.io/otel/trace v1.11.1
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
