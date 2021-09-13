package camunda

import "context"

const (
	VarTypeString  = "string"
	VarTypeDouble  = "double"
	VarTypeInteger = "integer"
)

type CamundaVariable struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type processStartParams struct {
	BusinessKey string                     `json:"businessKey"`
	Variables   map[string]CamundaVariable `json:"variables"`
}

type message struct {
	MessageName      string                     `json:"messageName"`
	BusinessKey      string                     `json:"businessKey"`
	ProcessVariables map[string]CamundaVariable `json:"processVariables"`
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
	WorkerID  string                     `json:"workerId"`
	Variables map[string]CamundaVariable `json:"variables"`
}

type Task struct {
	BusinessKey string                     `json:"business_key"`
	ID          string                     `json:"id"`
	TopicName   string                     `json:"topicName"`
	Variables   map[string]CamundaVariable `json:"variables"`
}

type (
	TaskCompleteFunc func(ctx context.Context, taskID string) error
	TaskHandlerFunc  func(ctx context.Context, completeFunc TaskCompleteFunc, t Task)
	Subscription     interface {
		Stop()
	}
)

func NewVariable(varType string, value interface{}) CamundaVariable {
	return CamundaVariable{
		Type:  varType,
		Value: value,
	}
}

func newMessage(messageName string, businessKey string, variables map[string]CamundaVariable) message {
	return message{
		MessageName:      messageName,
		BusinessKey:      businessKey,
		ProcessVariables: variables,
	}
}
