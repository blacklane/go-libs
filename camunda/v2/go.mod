module github.com/blacklane/go-libs/camunda/v2

go 1.16

replace (
	github.com/blacklane/go-libs/logger => ../../logger
	github.com/blacklane/go-libs/tracking => ../../tracking
	github.com/blacklane/go-libs/x/events => ../../x/events
)

require (
	github.com/blacklane/go-libs/logger v0.7.2
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.7.1
)

require github.com/stretchr/objx v0.1.1 // indirect
