GOLINT ?= golint
GO ?= go

.PHONY: test check deps
deps:
	$(GO) mod download

install_golint:
	$(GO) get -u golang.org/x/lint/golint

test: deps
	# musl build tag is required on Alpine linux to use the prebuilt librdkafka.
	# More on https://github.com/confluentinc/confluent-kafka-go/blob/master/README.md#librdkafka
	$(GO) test -race -tags musl -cover ./...

check:
	$(GO) vet $$($(GO) list ./...)
	$(GOLINT) $$($(GO) list ./...)
