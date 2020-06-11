GOLINT ?= golint
GO ?= go

.PHONY: test check deps
deps:
	$(GO) mod download

install_golint:
	$(GO) get -u golang.org/x/lint/golint

test: deps
	# using -musl because librdkafka is a C library
	$(GO) test -race -tags musl -cover ./...

check:
	$(GO) vet -tags musl $$($(GO) list musl ./...)
	$(GOLINT) $$($(GO) list musl ./...)
