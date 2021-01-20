package events

import (
	"context"
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-cmp/cmp"

	"github.com/blacklane/go-libs/tracking"
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

func TestAddTrackingID_EmptyContext(t *testing.T) {
	ctx := context.Background()
	event := Event{}

	addTrackingID(ctx, &event)

	if event.Headers != nil {
		t.Errorf("expected nil headers")
	}
}

func TestAddTrackingID_EmptyHeaders(t *testing.T) {
	wantID := "some-id"
	ctx := tracking.SetContextID(context.Background(), wantID)
	event := Event{}

	addTrackingID(ctx, &event)

	if event.Headers == nil {
		t.Errorf("expected headers to be set")
	}

	gotID := event.Headers[HeaderTrackingID]
	if !cmp.Equal(wantID, gotID) {
		t.Errorf("got key: %s, want: %s", gotID, wantID)
	}
}

func TestAddTrackingID_ExistingHeadersButNoTrackingID(t *testing.T) {
	wantID := "some-id"
	HeaderCustomID := "X-Custom-Id"
	customID := "custom-id"
	ctx := tracking.SetContextID(context.Background(), wantID)
	event := Event{
		Headers: Header{
			HeaderCustomID: customID,
		},
	}

	addTrackingID(ctx, &event)

	gotID := event.Headers[HeaderTrackingID]
	if !cmp.Equal(wantID, gotID) {
		t.Errorf("got key: %s, want: %s", gotID, wantID)
	}

	gotCustomID := event.Headers[HeaderCustomID]
	if !cmp.Equal(customID, gotCustomID) {
		t.Errorf("expected custom ID not to change. got key: %s, want: %s", gotCustomID, customID)
	}
}

func TestAddTrackingID_TrackingIDAlreadyExists(t *testing.T) {
	wantID := "some-id"
	otherID := "other-id"
	ctx := tracking.SetContextID(context.Background(), otherID)
	event := Event{
		Headers: Header{
			HeaderTrackingID: wantID,
		},
	}

	addTrackingID(ctx, &event)

	gotID := event.Headers[HeaderTrackingID]
	if cmp.Equal(otherID, gotID) {
		t.Errorf("expected tracking ID not to change. got key: %s, want: %s", gotID, otherID)
	}
	if !cmp.Equal(wantID, gotID) {
		t.Errorf("got key: %s, want: %s", gotID, wantID)
	}
}
