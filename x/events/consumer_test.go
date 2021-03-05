package events

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-cmp/cmp"
)

func TestParseHeaders(t *testing.T) {
	foo := "foo"
	bar := "bar"
	n42 := "42"
	s42 := "Answer to the Ultimate Question of Life, the Universe, and Everything"
	want := map[string]string{
		foo: bar,
		n42: s42,
	}

	got := parseHeaders([]kafka.Header{
		{Key: foo, Value: []byte(bar)},
		{Key: n42, Value: []byte(s42)},
	})

	if cmp.Equal(got, want) {
		t.Errorf("got: %v, want %v", got, want)
	}
}

func TestMessageToEvent(t *testing.T) {
	want := Event{
		Key: []byte("foo"),
		Headers: Header(map[string]string{
			"42": "Answer to the Ultimate Question of Life, the Universe, and Everything"}),
		Payload: []byte("bar"),
	}

	m := &kafka.Message{
		Value: []byte("bar"),
		Key:   []byte("foo"),
		Headers: []kafka.Header{{
			Key:   "42",
			Value: []byte("Answer to the Ultimate Question of Life, the Universe, and Everything")},
		},
	}

	got := messageToEvent(m)

	if cmp.Equal(got, want) {
		t.Errorf("got: %v, wnat %v", got, want)
	}

}

func TestKafkaConsumer_WithOAuth(t *testing.T) {
	kcc := NewKafkaConsumerConfig(&kafka.ConfigMap{"group.id": "TestKafkaConsumer_WithOAuth"})
	tokenSource := tokenSourceMock{}
	kcc.WithOAuth(tokenSource)

	c, err := NewKafkaConsumer(kcc, []string{"topic"})
	if err != nil {
		t.Fatalf("could not create kafka consumer: %v", err)
	}

	kc, ok := c.(*kafkaConsumer)
	if !ok {
		t.Fatalf("cannot cast Consumer to *kafkaConsumer")
	}

	if !cmp.Equal(kc.tokenSource, tokenSource) {
		t.Errorf("want: %v, got: %v", tokenSource, kc.tokenSource)
	}
}

func TestKafkaConsumer_WithErrFunc(t *testing.T) {
	var got error
	want := errors.New("error TestKafkaConfig_WithErrFunc")
	errFn := func(err error) { got = err }

	kcc := NewKafkaConsumerConfig(&kafka.ConfigMap{"group.id": "TestKafkaConsumer_WithErrFunc"})
	kcc.WithErrFunc(errFn)

	c, err := NewKafkaConsumer(kcc, []string{"topic"})
	if err != nil {
		t.Fatalf("cannot create kafka consumer: %v", err)
	}

	kc, ok := c.(*kafkaConsumer)
	if !ok {
		t.Fatalf("cannot cast Consumer to *kafkaConsumer")
	}

	kc.errFn(want)
	if got != want {
		t.Errorf("want: %v, got %v", want, got)
	}
}

func TestNewKafkaConsumerConfigAllInitialised(t *testing.T) {
	kc := NewKafkaConsumerConfig(nil)

	if kc.tokenSource == nil {
		t.Errorf("tokenSource is nil")
	}
	if kc.errFn == nil {
		t.Errorf("errFn is nil")
	}
}

func TestDeliverMessageOrderedPerKey(t *testing.T) {
	kcc := NewKafkaConsumerConfig(&kafka.ConfigMap{"group.id": "TestKafkaConsumer_WithErrFunc"})
	processedEvents := make([]string, 0)
	c, err := NewKafkaConsumer(kcc, []string{"topic"}, HandlerFunc(func(ctx context.Context, e Event) error {
		waitFor, _ := time.ParseDuration(string(e.Payload))
		time.Sleep(waitFor)
		processedEvents = append(processedEvents, fmt.Sprintf("%s:%s", e.Key, e.Payload))
		return nil
	}))
	if err != nil {
		t.Fatalf("cannot create kafka consumer: %v", err)
	}

	kc, ok := c.(*kafkaConsumer)
	if !ok {
		t.Fatalf("cannot cast Consumer to *kafkaConsumer")
	}

	msgWaiting3 := kafka.Message{
		Key:   []byte("key"),
		Value: []byte("1.5s"),
	}

	msgWaiting1 := kafka.Message{
		Key:   []byte("key"),
		Value: []byte("200ms"),
	}
	otherMsgWaiting1 := kafka.Message{
		Key:   []byte("other"),
		Value: []byte("1s"),
	}

	kc.deliverMessage(&msgWaiting3)
	kc.deliverMessage(&msgWaiting1)
	kc.deliverMessage(&otherMsgWaiting1)
	kc.wg.Wait()
	const expectedOrder = "other:1s,key:1.5s,key:200ms"
	outOrder := strings.Join(processedEvents, ",")
	if outOrder != expectedOrder {
		t.Errorf("incorrect order of processed messages %q when should be %q", outOrder, expectedOrder)
	}
}
