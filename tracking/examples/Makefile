KAFKA_GROUP_ID ?= $(shell echo "tracing-example-$$(date +%s%3)")
KAFKA_BOOTSTRAP_SERVERS ?=http://host.docker.internal:9092
KAFKA_TOPIC ?= tracing-example

GOPRIVATE ?= github.com/blacklane/*

deps: gitconfig
	GOPRIVATE=${GOPRIVATE} go mod download

gitconfig:
	@git config --global url.https://{GITHUB_TOKEN}:github.com.insteadOf https://github.com

.PHONY: run
run: deps
	go run *.go

#.PHONY: consumer-docker
#consumer-docker:
#	docker run --rm -it \
#		-e TOPIC=${KAFKA_TOPIC} \
#		-e BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} \
#		-e GROUP_ID=${KAFKA_GROUP_ID} \
# 		strimzi/kafka:0.14.0-kafka-2.3.0 \
# 		/opt/kafka/bin/kafka-console-consumer.sh \
#			--bootstrap-server ${KAFKA_BOOTSTRAP_SERVERS} \
#			--topic ${KAFKA_TOPIC} \
#			--from-beginning \
#			--consumer-property 'group.id=${KAFKA_GROUP_ID}$(date +%s%3)'
#
## When producing events they must follow the format MESSAGE_KEY:EVENT_MESSGE.
## depending on how the topic is configred the presence of a message key might
## be required.
#.PHONY: producer-docker
#producer-docker:
#	docker run --rm -it \
#		-e TOPIC=${KAFKA_TOPIC} \
#		-e BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} \
#		-e GROUP_ID=${KAFKA_GROUP_ID} \
# 		strimzi/kafka:0.14.0-kafka-2.3.0 \
#		/opt/kafka/bin/kafka-console-producer.sh \
#			--broker-list ${KAFKA_BOOTSTRAP_SERVERS} \
#			--topic ${KAFKA_TOPIC} \
#			--property 'parse.key=true' \
#			--property 'key.separator=:'

.PHONY: kafka-compose
kafka-compose:
	docker compose up \
    		--build \
    		--force-recreate \
    		--remove-orphans \
    		-d \
    		kafka

.PHONY: compose-down
compose-down clean:
	docker compose down

.PHONY: jaegertracing-compose
jaegertracing-compose:
	docker compose up \
    		--build \
    		--force-recreate \
    		--remove-orphans \
    		-d \
    		jaegertracing

.PHONY: compose-up
compose-up:
	docker compose up \
    		--build \
    		--force-recreate \
    		--remove-orphans \
    		-d

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

