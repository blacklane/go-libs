module github.com/blacklane/go-libs/logger

go 1.17

replace (
	github.com/blacklane/go-libs/tracking => ../tracking
	github.com/blacklane/go-libs/x/events => ../x/events
)

require (
	github.com/blacklane/go-libs/tracking v0.3.1
	github.com/blacklane/go-libs/x/events v0.5.0
	github.com/google/go-cmp v0.5.7
	github.com/rs/zerolog v1.29.0
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
)

require (
	github.com/confluentinc/confluent-kafka-go v1.9.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
