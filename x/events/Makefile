.PHONY: test check deps install_golint docker-tests docker-kafka clean

GOLINTER ?= revive
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
	$(GOLINTER) $$($(GO) list ./...)

docker-tests:
	export COMPOSE_PROJECT_NAME="golib-events-build-$(DRONE_BUILD_NUMBER)"
	docker compose --env-file .env_compose up \
		--build \
		--force-recreate \
		--remove-orphans \
		--abort-on-container-exit \
		--exit-code-from tests
	docker compose down

docker-kafka:
	docker compose --env-file .env_local up \
    		--build \
    		--force-recreate \
    		--remove-orphans \
    		-d \
    		kafka

docker-kafka-down clean:
	docker compose down
