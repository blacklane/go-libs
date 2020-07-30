package events

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/oauth2"
)

type emptyTokenSource struct{}

func (ts emptyTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

// kafkaCommon are common fields for KafkaConsumerConfig and KafkaProducerConfig.
type kafkaConfig struct {
	config      *kafka.ConfigMap
	tokenSource oauth2.TokenSource

	// errFn will be called for any error not handled by the handlers
	errFn func(error)
}

// kafkaCommon are common fields for Consumer and Producer.
type kafkaCommon struct {

	// errFn will be called for any error not handled by the handlers
	errFn func(error)

	tokenSource oauth2.TokenSource
}

// WithOAuth prepares to handle OAuth2.
// It'll set the kafka configurations:
//   sasl.mechanism: OAUTHBEARER
//   security.protocol: SASL_SSL
// it'll override any existing value for sasl.mechanism, security.protocol.
func (kc *kafkaConfig) WithOAuth(tokenSource oauth2.TokenSource) {
	kc.tokenSource = tokenSource
	setKafkaConfig(kc.config, "sasl.mechanism", "OAUTHBEARER")
	setKafkaConfig(kc.config, "security.protocol", "SASL_SSL")
}

// WithErrFunc sets a function to handle any error beyond producer delivery errors.
func (kc *kafkaConfig) WithErrFunc(errFn func(error)) {
	kc.errFn = errFn
}

func (kcp *kafkaCommon) refreshToken(handle kafka.Handle) {
	token, err := kcp.tokenSource.Token()
	if err != nil {
		errWrapped := fmt.Errorf("could not get oauth token: %w", err)

		kcp.errFn(errWrapped)
		err = handle.SetOAuthBearerTokenFailure(err.Error())
		if err != nil {
			kcp.errFn(fmt.Errorf("could not SetOAuthBearerTokenFailure: %w", errWrapped))
		}
	}

	err = handle.SetOAuthBearerToken(kafka.OAuthBearerToken{
		TokenValue: token.AccessToken,
		Expiration: token.Expiry,
	})
	if err != nil {
		kcp.errFn(fmt.Errorf("could not SetOAuthBearerToken: %w", err))
	}
}

func setKafkaConfig(configMap *kafka.ConfigMap, key string, val string) {
	_ = configMap.SetKey(key, val) // SetKey always return nil
}
