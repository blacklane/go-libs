package camunda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/blacklane/go-libs/logger"

	"github.com/blacklane/go-libs/x/camunda/internal"
)

type (
	Client interface {
		StartProcess(ctx context.Context, businessKey string, variables map[string]Variable) error
		SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]Variable)
		Subscribe(topicName string, handler TaskHandlerFunc, interval time.Duration) Subscription
	}

	camundaClient struct {
		camundaURL  string
		credentials *BasicAuthCredentials
		processKey  string
		httpClient  http.Client
		log         logger.Logger
	}
)

func NewClient(log logger.Logger, url string, processKey string, credentials *BasicAuthCredentials) Client {
	return &camundaClient{
		camundaURL:  url,
		credentials: credentials,
		processKey:  processKey,
		httpClient:  http.Client{},
		log:         log,
	}
}

func (c *camundaClient) StartProcess(
	ctx context.Context,
	businessKey string,
	variables map[string]Variable) error {
	if variables == nil {
		variables = map[string]Variable{}
	}

	variables[businessKeyVarKey] = Variable{
		Type:  VarTypeString,
		Value: businessKey,
	}
	params := processStartParams{
		BusinessKey: businessKey,
		Variables:   variables,
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(params)
	if err != nil {
		return fmt.Errorf("failed to send camunda message: could not encode parameters to json: %w",
			err)
	}

	url := fmt.Sprintf("process-definition/key/%s/start", c.processKey)
	_, err = c.doPostRequest(ctx, buf, url)
	if err != nil {
		c.log.Err(err).
			Str("business_key", params.BusinessKey).
			Msgf("failed to start process for business key - %s", params.BusinessKey)
		return err
	}
	c.log.Info().
		Str("business_key", params.BusinessKey).
		Msgf("process has been started with business key '%s'", params.BusinessKey)
	return nil
}

func (c *camundaClient) SendMessage(
	ctx context.Context,
	messageType string,
	businessKey string,
	updatedVariables map[string]Variable) {

	buf := new(bytes.Buffer)
	url := "message"
	newMessage := message{
		MessageName:      messageType,
		BusinessKey:      businessKey,
		ProcessVariables: updatedVariables,
	}

	err := json.NewEncoder(buf).Encode(newMessage)
	if err != nil {
		c.log.Err(err).
			Str(internal.LogFieldBusinessKey, businessKey).
			Msg("failed to send camunda message due to json error")
		return
	}

	_, err = c.doPostRequest(ctx, buf, url)
	if err != nil {
		c.log.Err(err).
			Str(internal.LogFieldBusinessKey, businessKey).
			Msgf("failed to send message for business key - %s", newMessage.BusinessKey)

		return
	}

	c.log.Info().
		Str(internal.LogFieldBusinessKey, businessKey).
		Msgf("%s message has been sent to camunda for business key '%s'", newMessage.MessageName, newMessage.BusinessKey)
}

func (c *camundaClient) Subscribe(
	topicName string,
	handler TaskHandlerFunc,
	interval time.Duration) Subscription {

	sub := &subscription{
		client:    c,
		topic:     topicName,
		isRunning: false,
		interval:  interval,
		log:       c.log,
	}
	sub.handler(handler)

	// run async fetch loop
	go sub.schedule()

	return sub
}

type ErrPost struct {
	Err        error
	Body       string
	HTTPStatus int
}

func (e ErrPost) Error() string {
	return e.Err.Error()
}

func (c *camundaClient) doPostRequest(ctx context.Context, params *bytes.Buffer, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.camundaURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, params)
	if err != nil {
		return nil, fmt.Errorf("could not create POST request to camunda engine: %w", err)
	}
	req.Header.Add(internal.HeaderContentType, "application/json")

	if c.credentials != nil {
		req.SetBasicAuth(c.credentials.User, c.credentials.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Err(err).Msgf("Could not send POST request")
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read reponse body: %w", err)
	}
	if resp.StatusCode >= 400 {
		// TODO: create a custom error type which carries the body so the user
		// can decide to inspect or log it.
		c.log.Error().
			Str(internal.LogFieldURL, url).
			Str("response_body", string(body)).
			Msgf("post request to camunda failed with status: %d", resp.StatusCode)
		return nil, ErrPost{
			Err:        fmt.Errorf("camunda API returned HTTP Status %d", resp.StatusCode),
			Body:       string(body),
			HTTPStatus: resp.StatusCode,
		}
	}

	return body, nil
}

func (c *camundaClient) complete(ctx context.Context, taskID string, params taskCompletionParams) error {
	log := logger.FromContext(ctx)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(params)
	if err != nil {
		log.Err(err).Msg("failed to complete camunda task due to json error")
		return err
	}

	url := fmt.Sprintf("external-task/%s/complete", taskID)
	_, err = c.doPostRequest(context.Background(), buf, url)
	if err != nil {
		return err
	}

	log.Debug().
		Str("task_id", taskID).
		Str("task_completion_params", fmt.Sprintf("%v", params)).
		Msgf("Completed task id %s", taskID)

	return nil
}

func (c *camundaClient) fetchAndLock(ctx context.Context, log logger.Logger, param fetchAndLock) ([]Task, error) {
	url := "external-task/fetchAndLock"

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(param)
	var tasks []Task
	if err != nil {
		return tasks, fmt.Errorf("failed to fetch camunda tasks: could not encode params: %w", err)
	}

	body, err := c.doPostRequest(ctx, buf, url)
	if err != nil {
		return tasks, err
	}

	err = json.Unmarshal(body, &tasks)
	if err != nil {
		log.Err(err).Msgf("Could not unmarshal task")
		return tasks, err
	}

	log.Debug().Msgf("Fetched %d Tasks from %s", len(tasks), c.camundaURL)
	return tasks, nil
}
