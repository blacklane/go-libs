package events

import (
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-cmp/cmp"
)

func TestKafkaProducer_WithOAuth(t *testing.T) {
	kpc := NewKafkaProducerConfig(&kafka.ConfigMap{})
	tokenSource := tokenSourceMock{}
	kpc.WithOAuth(tokenSource)

	p, err := NewKafkaProducer(kpc)
	if err != nil {
		t.Fatalf("could not create kafka consumer: %v", err)
	}

	kp, ok := p.(*kafkaProducer)
	if !ok {
		t.Fatalf("cannot cast Consumer to *kafkaProducer")
	}

	if !cmp.Equal(kp.tokenSource, tokenSource) {
		t.Errorf("want: %v, got: %v", tokenSource, kp.tokenSource)
	}
}

func TestKafkaProducer_WithErrFunc(t *testing.T) {
	var got error
	want := errors.New("error TestKafkaConfig_WithErrFunc")
	errFn := func(err error) { got = err }

	kpc := NewKafkaProducerConfig(&kafka.ConfigMap{"group.id": "TestKafkaConsumer_WithErrFunc"})
	kpc.WithErrFunc(errFn)

	p, err := NewKafkaProducer(kpc)
	if err != nil {
		t.Fatalf("cannot create kafka consumer: %v", err)
	}

	kp, ok := p.(*kafkaProducer)
	if !ok {
		t.Fatalf("cannot cast Consumer to *kafkaProducer")
	}

	kp.errFn(want)
	if got != want {
		t.Errorf("want: %v, got %v", want, got)
	}
}

func TestNewKafkaProducerConfigAllInitialised(t *testing.T) {
	kc := NewKafkaProducerConfig(nil)

	if kc.deliveryErrHandler == nil {
		t.Errorf("deliveryErrHandler is nil")
	}
	if kc.tokenSource == nil {
		t.Errorf("tokenSource is nil")
	}
	if kc.errFn == nil {
		t.Errorf("errFn is nil")
	}
}

func TestParseToKafkaHeaders(t *testing.T) {
	wantKey, wantValue := "HeaderName", "HeaderValue"
	eventHeaders := Header{
		wantKey: wantValue,
	}

	got := toKafkaHeaders(eventHeaders)[0]

	if !cmp.Equal(wantKey, got.Key) {
		t.Errorf("got key: %s, want: %s", got.Key, wantKey)
	}
	if !cmp.Equal(wantValue, string(got.Value)) {
		t.Errorf("got value: %s, want: %s", got.Value, wantValue)
	}
}
