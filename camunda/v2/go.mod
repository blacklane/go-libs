module github.com/blacklane/go-libs/camunda/v2

go 1.18

replace (
	github.com/blacklane/go-libs/logger => ../../logger
	github.com/blacklane/go-libs/tracking => ../../tracking
	github.com/blacklane/go-libs/x/events => ../../x/events
)

require (
	github.com/blacklane/go-libs/logger v0.7.2
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
)
