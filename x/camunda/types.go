package camunda

import "context"

const (
	VarTypeString  = "string"
	VarTypeDouble  = "double"
	VarTypeInteger = "integer"
)

type (
	Variable struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}

	processStartParams struct {
		BusinessKey string              `json:"businessKey"`
		Variables   map[string]Variable `json:"variables"`
	}

	message struct {
		MessageName      string              `json:"messageName"`
		BusinessKey      string              `json:"businessKey"`
		ProcessVariables map[string]Variable `json:"processVariables"`
	}

	topic struct {
		Name         string `json:"topicName"`
		LockDuration int    `json:"lockDuration"`
	}

	fetchAndLock struct {
		WorkerID string  `json:"workerId"`
		MaxTasks int     `json:"maxTasks"`
		Topics   []topic `json:"topics"`
	}

	taskCompletionParams struct {
		WorkerID  string              `json:"workerId"`
		Variables map[string]Variable `json:"variables"`
	}

	Task struct {
		BusinessKey string              `json:"business_key"`
		ID          string              `json:"id"`
		TopicName   string              `json:"topicName"`
		Variables   map[string]Variable `json:"variables"`
	}
)

type (
	TaskCompleteFunc func(ctx context.Context, taskID string) error
	TaskHandlerFunc  func(ctx context.Context, completeFunc TaskCompleteFunc, t Task)
	Subscription     interface {
		Stop()
	}
)
