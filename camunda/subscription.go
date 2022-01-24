package camunda

import (
	"context"
	"sync/atomic"
	"time"
)

type (
	Subscription struct {
		client           *client
		topic            string
		handlers         []TaskHandlerFunc
		interval         time.Duration
		isRunning        int32 // we need to use atomic int operations to make this thread safe
		maxTasksFetch    int
		taskLockDuration int
		workerID         string
	}
)

const (
	defaultFetchInterval    = time.Second * 5
	defaultMaxTasksFetch    = 100
	defaultTaskLockDuration = 2000 // 2seconds
	businessKeyJSONKey      = "BusinessKey"
)

func (s *Subscription) Stop() {
	atomic.StoreInt32(&s.isRunning, 0)
}

func MaxTasksFetch(maxTasksFetch int) func(subscription *Subscription) {
	return func(sub *Subscription) {
		sub.maxTasksFetch = maxTasksFetch
	}
}

func FetchInterval(interval time.Duration) func(subscription *Subscription) {
	return func(sub *Subscription) {
		sub.interval = interval
	}
}

func LockDuration(duration int) func(subscription *Subscription) {
	return func(sub *Subscription) {
		sub.taskLockDuration = duration
	}
}

func newSubscription(client *client, topic string, workerID string) *Subscription {
	return &Subscription{
		client:           client,
		topic:            topic,
		interval:         defaultFetchInterval,
		isRunning:        0,
		maxTasksFetch:    defaultMaxTasksFetch,
		taskLockDuration: defaultTaskLockDuration,
		workerID:         workerID,
	}
}

func (s *Subscription) complete(ctx context.Context, taskID string) error {
	completeParams := taskCompletionParams{
		WorkerID:  s.workerID,
		Variables: map[string]Variable{}, // we don't need to update any variables for now
	}

	return s.client.complete(ctx, taskID, completeParams)
}

// addHandler is attaching handlers to the Subscription
func (s *Subscription) addHandler(handler TaskHandler) {
	s.handlers = append(s.handlers, handler.Handle)
}

// fetch connects to camunda and starts polling the external tasks
// It will call each Handler if there is a new task on the topic
func (s *Subscription) fetch(ctx context.Context, fal fetchAndLock) {
	tasks, _ := s.client.fetchAndLock(&fal)
	for _, task := range tasks {
		for _, handler := range s.handlers {
			task.BusinessKey = extractBusinessKey(task)
			handler(ctx, s.complete, task)
		}
	}
}

func extractBusinessKey(task Task) string {
	value, ok := task.Variables[businessKeyJSONKey]
	if !ok {
		return ""
	}
	return value.Value.(string)
}

func (s *Subscription) schedule(ctx context.Context) {
	atomic.AddInt32(&s.isRunning, 1)
	lockParam := fetchAndLock{
		WorkerID: s.workerID,
		MaxTasks: s.maxTasksFetch,
		Topics: []topic{
			{
				Name:         s.topic,
				LockDuration: s.taskLockDuration,
			},
		},
	}

	for s.isRunning > 0 {
		<-time.After(s.interval) // fetch every x seconds
		go s.fetch(ctx, lockParam)
	}
}

func (h TaskHandlerFunc) Handle(ctx context.Context, completeFunc TaskCompleteFunc, t Task) {
	h(ctx, completeFunc, t)
}
