module github.com/blacklane/go-libs/logr/zerolog/example

go 1.18

require (
	github.com/blacklane/go-libs/logr v0.1.0
	github.com/blacklane/go-libs/logr/zerolog v0.1.0
	github.com/rs/zerolog v1.27.0
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20220614162138-6c1b26c55098 // indirect
)

replace (
	github.com/blacklane/go-libs/logr => ../../
	github.com/blacklane/go-libs/logr/zerolog => ../
)
