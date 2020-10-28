.PHONY: test check deps install_golint compose-tests compose-kafka clean

GOLINT ?= golint
GO ?= go

KAFKA_BOOTSTRAP_SERVERS ?= localhost:9092
# musl build tag is required on Alpine linux to use the prebuilt librdkafka.
# More on https://github.com/confluentinc/confluent-kafka-go/blob/master/README.md#librdkafka

BUILD_TAGS ?= -tags integration

deps:
	$(GO) mod download

install_golint:
	$(GO) get -u golang.org/x/lint/golint

test: deps
	KAFKA_BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} $(GO) test ${BUILD_TAGS} -cover ./...

check:
	$(GO) vet ${BUILD_TAGS} $$($(GO) list ./...)
	$(GOLINT) $$($(GO) list ./...)

compose-tests:
	docker-compose --env-file .env_compose up \
		--build \
		--force-recreate \
		--remove-orphans \
		--abort-on-container-exit \
		--exit-code-from tests
	docker-compose down

compose-kafka:
	docker-compose --env-file .env_local up \
    		--build \
    		--force-recreate \
    		--remove-orphans \
    		-d \
    		kafka

compose-kafka-down clean:
	docker-compose down
