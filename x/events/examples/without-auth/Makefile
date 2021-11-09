KAFKA_GROUP_ID ?=test-go-libs-events
KAFKA_BOOTSTRAP_SERVERS ?=localhost:9092
KAFKA_TOPIC ?=test-go-libs-events

.PHONY: run-producer
run-producer:
	KAFKA_BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} \
	KAFKA_GROUP_ID=${KAFKA_GROUP_ID} \
	KAFKA_TOPIC=${KAFKA_TOPIC} \
	go run producer/main.go

.PHONY: run-consumer
run-consumer:
	KAFKA_BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} \
	KAFKA_GROUP_ID=${KAFKA_GROUP_ID} \
	KAFKA_TOPIC=${KAFKA_TOPIC} \
	go run consumer/main.go

.PHONY: kafka-start
kafka-start:
	docker-compose --env-file ../../.env_local -f ../../docker-compose.yml up \
    		--build \
    		--force-recreate \
    		--remove-orphans \
    		-d \
    		kafka

.PHONY: kafka-stop
kafka-stop:
	docker-compose -f ../../docker-compose.yml down --remove-orphans --volumes
