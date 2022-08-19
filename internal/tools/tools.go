//go:build tools

package tools

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/jstemmer/go-junit-report/v2"
	_ "github.com/wadey/gocovmerge"
	_ "go.opentelemetry.io/build-tools/multimod"
	_ "golang.org/x/tools/cmd/stringer"
)
