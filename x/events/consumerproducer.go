package events

import (
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/oauth2"
)

// kafkaCP are common fields for KafkaConsumerConfig and KafkaProducerConfig.
type kafkaConfig struct {
	config      *kafka.ConfigMap
	tokenSource oauth2.TokenSource

	// errFn will be called for any error not handled by the handlers
	errFn func(error)
}

// kafkaCP are common fields for Consumer and Producer.
type kafkaCP struct {

	// errFn will be called for any error not handled by the handlers
	errFn func(error)

	tokenSource oauth2.TokenSource

	runMu    sync.RWMutex
	run      bool
	shutdown bool
}

func (kcp *kafkaCP) refreshToken(handle kafka.Handle) {

	token, err := kcp.tokenSource.Token()
	if err != nil {
		errWrapped := fmt.Errorf("could not get oauth token: %w", err)

		kcp.errFn(errWrapped)
		err = handle.SetOAuthBearerTokenFailure(err.Error())
		if err != nil {
			kcp.errFn(fmt.Errorf("could not SetOAuthBearerTokenFailure: %w", err))
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

func (kcp *kafkaCP) startRunning() bool {
	kcp.runMu.RLock()
	kcp.run = true
	defer kcp.runMu.RUnlock()
	return kcp.run
}

func (kcp *kafkaCP) running() bool {
	kcp.runMu.RLock()
	defer kcp.runMu.RUnlock()
	return kcp.run
}

func (kcp *kafkaCP) stopRun() {
	kcp.runMu.Lock()
	defer kcp.runMu.Unlock()
	kcp.run = false
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

// WithErrFunc sets a function to handle the kafka lib errors.
func (kc *kafkaConfig) WithErrFunc(errFn func(error)) {
	kc.errFn = errFn
}

func setKafkaConfig(configMap *kafka.ConfigMap, key string, val string) {
	_ = configMap.SetKey(key, val) // SetKey always return nil
}
