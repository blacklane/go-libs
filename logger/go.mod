module github.com/blacklane/go-libs/logger

go 1.17

replace (
	github.com/blacklane/go-libs/tracking => ../tracking
	github.com/blacklane/go-libs/x/events => ../x/events
)

require (
	github.com/blacklane/go-libs/tracking v0.3.1
	github.com/blacklane/go-libs/x/events v0.2.1
	github.com/google/go-cmp v0.5.7
	github.com/rs/zerolog v1.26.0
	github.com/stretchr/testify v1.7.1
)

require (
	github.com/confluentinc/confluent-kafka-go v1.8.2 // indirect
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	golang.org/x/net v0.0.0-20211109214657-ef0fda0de508 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
