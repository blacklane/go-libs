module github.com/blacklane/go-libs/logr/zerolog/example

go 1.18

require (
	github.com/blacklane/go-libs/logr v0.1.0
	github.com/blacklane/go-libs/logr/zap v0.1.0
	go.uber.org/zap v1.21.0
)

require (
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
)

replace (
	github.com/blacklane/go-libs/logr => ../../
	github.com/blacklane/go-libs/logr/zap => ../
)
