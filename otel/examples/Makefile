KAFKA_GROUP_ID ?= $(shell echo "tracing-example-$$(date +%s%3)")
KAFKA_BOOTSTRAP_SERVERS ?=http://host.docker.internal:9092
KAFKA_TOPIC ?= tracing-example

GOPRIVATE ?= github.com/blacklane/*
GOLINTER ?= revive

BUILD_TAGS ?=
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	BUILD_TAGS =-tags musl,integration
endif

.PHONY: deps
deps: gitconfig
	GOPRIVATE=${GOPRIVATE} go mod download

.PHONY: gitconfig
gitconfig:
	@git config --global url.https://${GITHUB_TOKEN}@github.com.insteadOf https://github.com

.PHONY: run-http
run-http: deps
	go run ${BUILD_TAGS} cmd/http/main.go

.PHONY: run-events
run-events: deps
	go run ${BUILD_TAGS} cmd/events/main.go

.PHONY: lint
lint:
	${GOLINTER} ./...

.PHONY: compose-dependencies
compose-dependencies:
	docker-compose up \
    		--remove-orphans \
    		-d \
    		elasticsearch kafka zookeeper jaeger

.PHONY: compose-down
compose-down:
	docker-compose down

.PHONY: clean
clean:
	docker-compose down --volumes --remove-orphans

.PHONY: jaegertracing-docker
jaegertracing-docker:
	docker run -d \
	  --name jaeger \
	  --rm -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
	  -p 5775:5775/udp \
	  -p 6831:6831/udp \
	  -p 6832:6832/udp \
	  -p 5778:5778 \
	  -p 16686:16686 \
	  -p 14268:14268 \
	  -p 14250:14250 \
	  -p 9411:9411 \
	  jaegertracing/all-in-one:1.21
