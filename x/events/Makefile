GOPRIVATE ?=
GOLINT ?= golint
GO ?= go

deps:
	GOPRIVATE=${GOPRIVATE} $(GO) mod download

install_golint:
	$(GO) get -u golang.org/x/lint/golint

.PHONY: test
test: deps
	$(GO) test -race -tags musl -cover ./...

.PHONY: check
check:
	$(GO) vet -tags musl $$($(GO) list -tags musl ./...)
	$(GOLINT) $$($(GO) list -tags musl ./...)
