package camunda

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testTopic = "test-topic"
	workerID  = "test-worker"
)

func TestSubscription_Stop(t *testing.T) {
	client, _ := newTestClient()
	sub := newSubscription(client, testTopic, workerID)

	// fake start fetching
	sub.isRunning = 1

	// act
	sub.Stop()

	assert.Equal(t, int32(0), sub.isRunning)
}

func TestFetchInterval(t *testing.T) {
	client, _ := newTestClient()
	sub := newSubscription(client, testTopic, workerID)

	// act
	FetchInterval(time.Minute)(sub)

	assert.Equal(t, time.Minute, sub.interval)
}

func TestLockDuration(t *testing.T) {
	client, _ := newTestClient()
	sub := newSubscription(client, testTopic, workerID)
	newLockDuration := 42

	// act
	LockDuration(newLockDuration)(sub)

	assert.Equal(t, newLockDuration, sub.taskLockDuration)
}

func TestMaxTasksFetch(t *testing.T) {
	client, _ := newTestClient()
	sub := newSubscription(client, testTopic, workerID)
	newMaxTasksFetch := 42

	// act
	MaxTasksFetch(newMaxTasksFetch)(sub)

	assert.Equal(t, newMaxTasksFetch, sub.maxTasksFetch)
}

func TestSubscription_complete_handlesError(t *testing.T) {
	client, mockHttpClient := newTestClient()
	sub := newSubscription(client, testTopic, workerID)
	taskID := uuid.New().String()

	body := io.NopCloser(bytes.NewReader([]byte("{}")))
	mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 404,
		Body:       body,
	})

	// act
	err := sub.complete(context.Background(), taskID)

	assert.Error(t, err)
}
