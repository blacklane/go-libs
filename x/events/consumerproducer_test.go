package events

import (
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
)

func TestKafkaConfig_WithErrFunc(t *testing.T) {
	var got error
	want := errors.New("error TestKafkaConfig_WithErrFunc")
	errFn := func(err error) { got = err }

	kc := kafkaConfig{}

	kc.WithErrFunc(errFn)

	kc.errFn(want)
	if got != want {
		t.Errorf("want: %v, got %v", want, got)
	}
}

type tokenSourceMock struct{}

func (mt tokenSourceMock) Token() (*oauth2.Token, error) { return nil, nil }

func TestKafkaConfig_WithOAuth(t *testing.T) {
	assert := func(t *testing.T, kc kafkaConfig, key string, want string) {
		missing := "missing"
		val, err := kc.config.Get(key, missing)
		if err != nil {
			t.Errorf("error getting config '%s': %v", key, err)
		}
		s, ok := val.(string)
		if !ok {
			t.Errorf("cannot cast kafka config '%s':%v to string", key, val)
		}
		if s != want {
			t.Errorf("kafka config %s: want: %s, got: %s", key, want, s)
		}
	}

	tokenSource := tokenSourceMock{}
	kc := kafkaConfig{
		config: &kafka.ConfigMap{},
	}

	kc.WithOAuth(tokenSource)

	assert(t, kc, "sasl.mechanism", "OAUTHBEARER")
	assert(t, kc, "security.protocol", "SASL_SSL")

	if !cmp.Equal(kc.tokenSource, tokenSource) {
		t.Errorf("want: %v, got: %v", tokenSource, kc.tokenSource)
	}
}
