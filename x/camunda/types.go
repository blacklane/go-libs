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
	Task struct {
		BusinessKey string              `json:"business_key"`
		ID          string              `json:"id"`
		TopicName   string              `json:"topicName"`
		Variables   map[string]Variable `json:"variables"`
	}
	TaskHandler interface {
		Handle(completeFunc TaskCompleteFunc, t Task)
	}
	TaskCompleteFunc func(ctx context.Context, taskID string) error
	TaskHandlerFunc  func(completeFunc TaskCompleteFunc, t Task)
)

func NewVariable(varType string, value interface{}) Variable {
	return Variable{
		Type:  varType,
		Value: value,
	}
}

type processStartParams struct {
	BusinessKey string              `json:"businessKey"`
	Variables   map[string]Variable `json:"variables"`
}

type message struct {
	MessageName      string              `json:"messageName"`
	BusinessKey      string              `json:"businessKey"`
	ProcessVariables map[string]Variable `json:"processVariables"`
}

type topic struct {
	Name         string `json:"topicName"`
	LockDuration int    `json:"lockDuration"`
}

type fetchAndLock struct {
	WorkerID string  `json:"workerId"`
	MaxTasks int     `json:"maxTasks"`
	Topics   []topic `json:"topics"`
}

type taskCompletionParams struct {
	WorkerID  string              `json:"workerId"`
	Variables map[string]Variable `json:"variables"`
}

func newMessage(messageName string, businessKey string, variables map[string]Variable) message {
	return message{
		MessageName:      messageName,
		BusinessKey:      businessKey,
		ProcessVariables: variables,
	}
}
