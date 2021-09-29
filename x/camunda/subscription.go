package camunda

import (
	"context"
	"fmt"
	"time"

	"github.com/blacklane/go-libs/logger"
	"github.com/blacklane/go-libs/tracking"
	"github.com/rs/zerolog"

	"github.com/blacklane/go-libs/x/camunda/internal"
)

type (
	// is there any good reason to implement the methods with a pointer receiver?
	subscription struct {
		client    *client
		topic     string
		handlers  []TaskHandlerFunc
		isRunning bool
		interval  time.Duration
		log       logger.Logger
	}
)

const (
	// TODO: accept is as parameter in the constructor
	workerID         = "ride-bundler-worker"
	maxTasksFetch    = 100
	taskLockDuration = 10000 // 10s
	// TODO: make it configurable
	// I'm confused, what KeyVarKey means?
	businessKeyVarKey = "BusinessKey"
)

// shouldn't it receive the variables as well?
func (s *subscription) complete(ctx context.Context, taskID string) error {
	completeParams := taskCompletionParams{
		WorkerID:  workerID,
		Variables: map[string]CamundaVariable{}, // we don't need to update any variables for now
	}

	return s.client.complete(ctx, taskID, completeParams)
}

// this should probably be called 'addHandler' as it is not a handler and calling
// it handler is either confusing or misleading.
// Handler is attaching handlers to the Subscription
func (s *subscription) handler(handler TaskHandlerFunc) {
	s.handlers = append(s.handlers, handler)
}

// please follow Go's standard for doc comments:
//   Doc comments work best as complete sentences, which allow a wide variety of
//   automated presentations. The first sentence should be a one-sentence summary
//   that starts with the name being declared.
//   https://golang.org/doc/effective_go#commentary
// besides, did u rename this method? It's a bit off the doc comment

// Open connects to camunda and start polling the external tasks
// It will call each handler if there is a new task on the topic
func (s *subscription) fetch(fal fetchAndLock) {
	tasks, _ := s.client.fetchAndLock(s.log, &fal)
	for _, task := range tasks {
		for _, handler := range s.handlers {
			ctx := s.getContextForTask(task)
			task.BusinessKey = extractBusinessKey(task)
			// better to run in a goroutine, no need to block for every handler
			handler(ctx, s.complete, task)
		}
	}
}

// please follow Go's standard for doc comments:
//   Doc comments work best as complete sentences, which allow a wide variety of
//   automated presentations. The first sentence should be a one-sentence summary
//   that starts with the name being declared.
//   https://golang.org/doc/effective_go#commentary

// create context with logger & tracking ID
func (s *subscription) getContextForTask(task Task) context.Context {
	trackingID := fmt.Sprintf("camunda-task-%s-%s", s.topic, task.ID)
	ctx := tracking.SetContextID(context.Background(), trackingID)

	newLogger := s.log.With().Logger()
	newLogger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(internal.LogFieldCamundaTaskID, task.ID).
			Str(internal.LogFieldBusinessKey, task.BusinessKey)
	})
	return newLogger.WithContext(ctx)
}

func extractBusinessKey(task Task) string {
	// var empty CamundaVariable is confusing, just compare against CamundaVariable{}
	// if task.Variables[businessKeyVarKey] == CamundaVariable{}
	//
	// or even better, use the “comma ok” idiom:
	// 	val, ok := task.Variables[businessKeyVarKey]
	//	if !ok {
	//		return ""
	//	}
	//
	//	return val.Value.(string)
	var empty CamundaVariable
	if task.Variables[businessKeyVarKey] == empty {
		return ""
	}
	return task.Variables[businessKeyVarKey].Value.(string)
}

func (s *subscription) schedule() {
	s.isRunning = true
	lockParam := fetchAndLock{
		WorkerID: workerID,
		MaxTasks: maxTasksFetch,
		Topics: []topic{
			{
				Name:         s.topic,
				LockDuration: taskLockDuration,
			},
		},
	}

	// Blocker: there is a race condition here
	for s.isRunning {
		<-time.After(s.interval) // fetch every x seconds
		go s.fetch(lockParam)
	}
}

func (s *subscription) Stop() {
	s.isRunning = false
}
