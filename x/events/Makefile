GOPRIVATE ?=
GOLINT ?= golint
GO ?= go

deps:
	GOPRIVATE=${GOPRIVATE} $(GO) mod download

install_golint:
	$(GO) get -u golang.org/x/lint/golint

test: deps
	$(GO) test -tags musl -cover ./...

check:
	$(GO) vet -tags musl $$($(GO) list -tags musl ./...)
	$(GOLINT) $$($(GO) list -tags musl ./...)
