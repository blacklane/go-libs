package camunda

import (
	"context"
)

const (
	VarTypeString  = "string"
	VarTypeDouble  = "double"
	VarTypeInteger = "integer"
)

// A general suggestion: move all exported names to the top, usually something like:
//   consts
//   vars
//   interfaces
//   methods

// Rename it to variable, this name is actually `camunda.CamundaVariable`,
// what is redundant. Besides:
//    The importer of a package will use the name to refer to its contents, so
//    exported names in the package can use that fact to avoid repetition.
//    (Don't use the import . notation, which can simplify tests that must run
//    outside the package they are testing, but should otherwise be avoided.) For
//    instance, the buffered reader type in the bufio package is called Reader, not
//    BufReader, because users see it as bufio.Reader, which is a clear, concise name.
//    https://golang.org/doc/effective_go#package-names

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
	// I find it a bit odd this approach for handling tasks using
	// closures. And receiving a function to complete it. I did not see the value on it.
	// I'd make the task handler an interface. I'll give freedom to the user
	// to implement it the way that is more convenient for them. And I'd provide
	// a handlerFunc adapter just as the http.HandlerFunc type. It'll let the
	// user free to use whatever is easier for them.
	// For the complete function, I'm not sure, I believe ethe better would be to
	// discuss hw it can be used. However, from what I can see here I'd try to make
	// it either a method on Task (which might be complicated because of the need
	// of a http.Client) or maybe a channel and another goroutine takes care of
	// completing the task. There is still the error to handle though. As I said
	// we could talk more about this and brainstorm other design options.

	TaskCompleteFunc func(ctx context.Context, taskID string) error
	TaskHandlerFunc  func(ctx context.Context, completeFunc TaskCompleteFunc, t Task)
	Subscription     interface {
		Stop()
	}
)

// Please, remove it. There s no need for that, this function does nothing. It's
// even more clear to instantiate the type directly in the code. If I see a `NewSomething`
// I cannot be sure what it returns, a concrete type, an interface, a concrete type
// which implements an interface, does it validate something? There are many possibilities
// for a `New...` function, whereas constructing the type directly is 100% explicit.
// Another thing, do you see how `CamundaVariable` should actually be `Variable`?
// You didn't name the function `NewCamundaVariable`

func NewVariable(varType string, value interface{}) CamundaVariable {
	return CamundaVariable{
		Type:  varType,
		Value: value,
	}
}

// same as above, besides it's a not exported method, no need at all for it.
func newMessage(messageName string, businessKey string, variables map[string]CamundaVariable) message {
	return message{
		MessageName:      messageName,
		BusinessKey:      businessKey,
		ProcessVariables: variables,
	}
}
