package camunda

import (
	"time"
)

type (
	subscription struct {
		client    *client
		topic     string
		handlers  []TaskHandlerFunc
		isRunning bool
		interval  time.Duration
	}
)

const (
	workerID          = "ride-bundler-worker"
	maxTasksFetch     = 100
	taskLockDuration  = 10000 // 10s
	businessKeyVarKey = "BusinessKey"
)

func (s *subscription) complete(taskID string) error {
	completeParams := taskCompletionParams{
		WorkerID:  workerID,
		Variables: map[string]CamundaVariable{}, // we don't need to update any variables for now
	}

	return s.client.complete(taskID, completeParams)
}

// Handler is attaching handlers to the Subscription
func (s *subscription) handler(handler TaskHandlerFunc) {
	s.handlers = append(s.handlers, handler)
}

// Open connects to camunda and start polling the external tasks
// It will call each handler if there is a new task on the topic
func (s *subscription) fetch(fal fetchAndLock) {
	tasks, _ := s.client.fetchAndLock(&fal)
	for _, task := range tasks {
		for _, handler := range s.handlers {
			task.BusinessKey = extractBusinessKey(task)
			handler(s.complete, task)
		}
	}
}

func extractBusinessKey(task Task) string {
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

	for s.isRunning {
		<-time.After(s.interval) // fetch every x seconds
		go s.fetch(lockParam)
	}
}

func (s *subscription) Stop() {
	s.isRunning = false
}
