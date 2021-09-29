package camunda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/blacklane/go-libs/logger"

	"github.com/blacklane/go-libs/x/camunda/internal"
)

type (
	Client interface {
		StartProcess(ctx context.Context, businessKey string, variables map[string]CamundaVariable) error
		SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]CamundaVariable)
		Subscribe(topicName string, handler TaskHandlerFunc, interval time.Duration) Subscription
	}

	// you can call it `client`, prefixing with 'camunda' is redundant
	client struct {
		camundaURL  string
		credentials *BasicAuthCredentials
		processKey  string
		httpClient  HttpClient
		log         logger.Logger
	}
)

// I would not receive a logger here. First, as a rule of thumb, a lib shouldn't
// be logging, better to return a good error and let the user decide what to do.
// Second, the logger does not have "context", it does not have trackingID, and so on.
// If it's really needed to log, better to get it from the context.

// Blocker: it needs to receive a http.Client
// Why is credentials a pointer to BasicAuthCredentials?
// For configuration, mainly optional ones, I recommend to use functional options.
// The last parameter is a variadic parameter of a function that applies the options.
// I've already resented this pattern at least once, and here is a good blog post
// about it: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
// and a video: https://www.youtube.com/watch?v=eyVDTmB9tpE

func NewClient(log logger.Logger, url string, processKey string, credentials *BasicAuthCredentials) Client {
	return &client{
		camundaURL:  url,
		credentials: credentials,
		processKey:  processKey,
		httpClient:  &http.Client{},
		log:         log,
	}
}

func (c *client) StartProcess(ctx context.Context, businessKey string, variables map[string]CamundaVariable) error {
	variables[businessKeyVarKey] = NewVariable(VarTypeString, businessKey)
	params := processStartParams{
		BusinessKey: businessKey,
		Variables:   variables,
	}

	// this can be simplified to:
	// 	buf := bytes.Buffer{}
	//	err := json.NewEncoder(&buf).Encode(params)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(params)
	if err != nil {
		// please do not log, wrapping the error is better, let the user handle
		// error. Besides, it isn't a "json error", it's a parsing error.
		// return fmt.Errorf("failed to send camunda message, could not parse eparams to JSON: %w",err)
		c.log.Err(err).Msg("failed to send camunda message due to json error")
		return err
	}

	url := fmt.Sprintf("process-definition/key/%s/start", c.processKey)
	_, err = c.doPostRequest(ctx, buf, url)
	if err != nil {
		// please do not log, wrapping the error is better
		c.log.Err(err).
			Msgf("failed to start process for business key - %s", params.BusinessKey)
		return err
	}
	// please do not log, it's a debug log, let the user decide if they need it
	// or not. Either take a debug configuration or log as debug
	c.log.Info().Msgf("process has been started with business key '%s'", params.BusinessKey)
	return nil
}

func (c *client) SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]CamundaVariable) {
	buf := new(bytes.Buffer)
	url := "message"
	newMessage := newMessage(messageType, businessKey, updatedVariables)
	err := json.NewEncoder(buf).Encode(newMessage)
	if err != nil {
		c.log.Err(err).
			Str(internal.LogFieldBusinessKey, businessKey).
			Msg("failed to send camunda message due to json error")
		// do not log and return the error wrapped on this log message.
		// Another thing, it's a parsing error, not a json error
		return
	}

	_, err = c.doPostRequest(ctx, buf, url)
	if err != nil {
		c.log.Err(err).
			Str(internal.LogFieldBusinessKey, businessKey).
			Msgf("failed to send message for business key - %s", newMessage.BusinessKey)
		// do not log and return the error wrapped on this log message.
		// Another thing, it's a parsing error, not a json error
		return
	}

	c.log.Info().
		Str(internal.LogFieldBusinessKey, businessKey).
		Msgf("%s message has been sent to camunda for business key '%s'", newMessage.MessageName, newMessage.BusinessKey)
}

func (c *client) Subscribe(topicName string, handler TaskHandlerFunc, interval time.Duration) Subscription {
	// this is more a debug log, remove it. better to receive a debug option and only then log it.
	// if the user wants to log it, they can do it just before calling Subscribe.
	c.log.Info().Msgf("Subscribing to: %s", topicName)
	sub := &subscription{
		client:    c,
		topic:     topicName,
		isRunning: false,
		interval:  interval,
		log:       c.log,
	}
	// see how confusing it is? a handler receiving a handler. I cannot even try
	// to infer what it does, I'll need to look at the implementation. However,
	// if it were called `addHandler`, I'd know sub.addHandler(handler) is adding
	// a handler
	sub.handler(handler)

	// run async fetch loop
	go sub.schedule()

	return sub
}

func (c *client) doPostRequest(ctx context.Context, params *bytes.Buffer, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.camundaURL, endpoint)

	// use the constant http.MethodPost
	// better to use http.NewRequestWithContext
	req, err := http.NewRequest("POST", url, params)
	if err != nil {
		c.log.Err(err).Msgf("Could not create POST request")
		return nil, err
	}
	req.Header.Add(internal.HeaderContentType, "application/json")

	if c.credentials != nil {
		req.SetBasicAuth(c.credentials.User, c.credentials.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// do not log, let the user handle the error, the user should not decide how to handle it.
		// Besides, just imagine, as a user of this lib you have got no idea the lib
		// logs this message. Then suddenly your served starts to log a generic error "Could not send POST request".
		// How can the user know it comes from the lib, they'd need to read the code up to here to find it out.
		c.log.Err(err).Msgf("Could not send POST request")
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		// blocker: do not log the body. As it's useful, better to return it on
		// a custom error type
		// please do not log, return a wrapped error:
		// return nil, fmt.Errorf("camunda API returned HTTP status %d", resp.StatusCode)
		// you can have a custom error like:
		// type ErrHttp struct {
		//	Err        error
		//	Body       string
		//	StatusCode int
		// }
		//
		// func (e ErrHttp) Error() string {
		//	return fmt.Sprintf("http status code: %d, %s", e.StatusCode, e.Err.Error())
		// }
		// return nil, ErrHttp{
		//			Err:        fmt.Errorf("camunda API returned HTTP status %d", resp.StatusCode),
		//			Body:       string(body),
		//			StatusCode: resp.StatusCode,
		//		}
		c.log.Error().Str(internal.LogFieldURL, url).Msgf("Status: %v, Body: %v", resp.StatusCode, string(body))
		// use fmt.Errorf instead of errors.New(fmt.Sprintf(...
		return nil, errors.New(fmt.Sprintf("Camunda API returned Status %v", resp.StatusCode))
	}

	return body, nil
}

func (c *client) complete(ctx context.Context, taskId string, params taskCompletionParams) error {
	log := logger.FromContext(ctx)

	// this can be simplified to:
	// 	buf := bytes.Buffer{}
	//	err := json.NewEncoder(&buf).Encode(params)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(params)
	if err != nil {
		// do not log, return a wrapped error:
		// return fmt.Errorf("could not parse taskCompletionParamss to JSON: %w", err)
		log.Err(err).Msg("failed to complete camunda task due to json error")
		return err
	}

	url := fmt.Sprintf("external-task/%s/complete", taskId)
	_, err = c.doPostRequest(context.Background(), buf, url)
	if err != nil {
		return err
	}

	// remove this log, it's a debug log. If you want to provide this functionality
	// add a debug option when creating the client or log it as debug, the former is better
	log.Info().Msgf("Completed task: %v", params.Variables)

	return nil
}

func (c *client) fetchAndLock(log logger.Logger, param *fetchAndLock) ([]Task, error) {
	// this can be simplified to:
	// 	buf := bytes.Buffer{}
	//	err := json.NewEncoder(&buf).Encode(params)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(param)
	var tasks []Task // it can be declared just above json.Unmarshal. It's better as it isn't used before that
	if err != nil {
		// same as the others, return a wrapped error
		log.Err(err).Msg("failed to fetch camunda tasks due to json error")
		// it's and error case, return nil, error. I know tasks is nil
		// at this point of the code, but first it's more explicit to return nil,
		// second, if the code changes, it might introduce a bug or unintended
		// behaviour where the returned tasks isn't nil and should be
		return tasks, err
	}

	url := "external-task/fetchAndLock"
	body, err := c.doPostRequest(context.Background(), buf, url)
	if err != nil {
		// same as the others, return nil, wrapped error
		return tasks, err
	}

	err = json.Unmarshal(body, &tasks)
	if err != nil {
		// same as the others, return nil, wrapped error
		log.Err(err).Msgf("Could not unmarshal task")
		return tasks, err
	}

	// remove this log, it's a debug log. If you want to provide this functionality
	// add a debug option when creating the client or log it as debug, the former is better
	log.Debug().Msgf("Fetched %d Tasks from %s", len(tasks), c.camundaURL)
	return tasks, nil
}
