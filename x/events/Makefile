.PHONY: test check deps install_golint

GOLINT ?= golint
GO ?= go

deps:
	$(GO) mod download

install_golint:
	$(GO) get -u golang.org/x/lint/golint

test: deps
	# musl build tag is required on Alpine linux to use the prebuilt librdkafka.
	# More on https://github.com/confluentinc/confluent-kafka-go/blob/master/README.md#librdkafka
	$(GO) test -race -cover ./...

check:
	$(GO) vet -tags musl $$($(GO) list ./...)
	$(GOLINT) $$($(GO) list ./...)
