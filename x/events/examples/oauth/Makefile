.PHONY: consumer-docker producer-docker

OAUTH_CLIENT_ID ?=test-go-libs-events
OAUTH_CLIENT_SECRET ?=FIX_ME_CLIENT_ID
# TODO(Anderson): remove this default
OAUTH_TOKEN_ENDPOINT_URI ?=https://keycloak.int.blacklane.io/auth/realms/master/protocol/openid-connect/token
KAFKA_GROUP_ID ?=test-go-libs-events
# TODO(Anderson): remove this default
KAFKA_BOOTSTRAP_SERVERS ?=kafka-bootstrap.int.blacklane.io:9094
KAFKA_TOPIC ?=test-go-libs-events

run-producer:
	@KAFKA_BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} \
	KAFKA_GROUP_ID=${KAFKA_GROUP_ID} \
	KAFKA_TOPIC=${KAFKA_TOPIC} \
	OAUTH_CLIENT_ID=${OAUTH_CLIENT_ID} \
	OAUTH_CLIENT_SECRET=${OAUTH_CLIENT_SECRET} \
	OAUTH_TOKEN_ENDPOINT_URI=${OAUTH_TOKEN_ENDPOINT_URI} \
	go run producer/main.go

run-consumer:
	@KAFKA_BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} \
	KAFKA_GROUP_ID=${KAFKA_GROUP_ID} \
	KAFKA_TOPIC=${KAFKA_TOPIC} \
	OAUTH_CLIENT_ID=${OAUTH_CLIENT_ID} \
	OAUTH_CLIENT_SECRET=${OAUTH_CLIENT_SECRET} \
	OAUTH_TOKEN_ENDPOINT_URI=${OAUTH_TOKEN_ENDPOINT_URI} \
	go run consumer/main.go

kafka:
	docker compose --env-file ../../.env_local -f ../../docker-compose.yml up \
    		--build \
    		--force-recreate \
    		--remove-orphans \
    		-d \
    		kafka

consumer-docker:
	docker run --rm -it \
		-e CLASSPATH='/opt/kafka/libs/kafka-oauth-client-*:/opt/kafka/libs/kafka-oauth-common-*' \
		-e OAUTH_CLIENT_ID=${OAUTH_CLIENT_ID} \
		-e OAUTH_CLIENT_SECRET=${OAUTH_CLIENT_SECRET} \
		-e OAUTH_TOKEN_ENDPOINT_URI=${OAUTH_TOKEN_ENDPOINT_URI} \
		-e TOPIC=${KAFKA_TOPIC} \
		-e BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} \
		-e GROUP_ID=${KAFKA_GROUP_ID} \
 		strimzi/kafka:0.14.0-kafka-2.3.0 \
 		/opt/kafka/bin/kafka-console-consumer.sh \
			--bootstrap-server ${KAFKA_BOOTSTRAP_SERVERS} \
			--topic ${KAFKA_TOPIC} \
			--from-beginning \
			--consumer-property 'group.id=${KAFKA_GROUP_ID}' \
			--consumer-property 'security.protocol=SASL_SSL' \
			--consumer-property 'sasl.mechanism=OAUTHBEARER' \
			--consumer-property 'sasl.jaas.config=org.apache.kafka.common.security.oauthbearer.OAuthBearerLoginModule required;' \
			--consumer-property 'sasl.login.callback.handler.class=io.strimzi.kafka.oauth.client.JaasClientOauthLoginCallbackHandler'

# When producing events they must follow the format MESSAGE_KEY:EVENT_MESSAGE.
# Depending on how the topic is configured the presence of a message key is optional.
producer-docker:
	docker run --rm -it \
		-e CLASSPATH='/opt/kafka/libs/kafka-oauth-client-*:/opt/kafka/libs/kafka-oauth-common-*' \
		-e OAUTH_CLIENT_ID=${OAUTH_CLIENT_ID} \
		-e OAUTH_CLIENT_SECRET=${OAUTH_CLIENT_SECRET} \
		-e OAUTH_TOKEN_ENDPOINT_URI=${OAUTH_TOKEN_ENDPOINT_URI} \
		-e TOPIC=${KAFKA_TOPIC} \
		-e BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS} \
		-e GROUP_ID=${KAFKA_GROUP_ID} \
 		strimzi/kafka:0.14.0-kafka-2.3.0 \
		/opt/kafka/bin/kafka-console-producer.sh \
			--broker-list ${KAFKA_BOOTSTRAP_SERVERS} \
			--topic ${KAFKA_TOPIC} \
			--property "parse.key=true" \
			--property "key.separator=:" \
			--producer-property 'security.protocol=SASL_SSL' \
			--producer-property 'sasl.mechanism=OAUTHBEARER' \
			--producer-property 'sasl.jaas.config=org.apache.kafka.common.security.oauthbearer.OAuthBearerLoginModule required;' \
			--producer-property 'sasl.login.callback.handler.class=io.strimzi.kafka.oauth.client.JaasClientOauthLoginCallbackHandler'
