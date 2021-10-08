package camunda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/blacklane/go-libs/x/camunda/internal"
)

type (
	Client interface {
		StartProcess(ctx context.Context, businessKey string, variables map[string]CamundaVariable) error
		SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]CamundaVariable) error
		Subscribe(topicName string, handler TaskHandlerFunc, interval time.Duration) Subscription
	}
	client struct {
		camundaURL  string
		credentials *internal.BasicAuthCredentials
		processKey  string
		httpClient  internal.HttpClient
	}
)

func NewClient(url string, processKey string, credentials *internal.BasicAuthCredentials) Client {
	return &client{
		camundaURL:  url,
		credentials: credentials,
		processKey:  processKey,
		httpClient:  &http.Client{},
	}
}

func (c *client) StartProcess(ctx context.Context, businessKey string, variables map[string]CamundaVariable) error {
	variables[businessKeyVarKey] = NewVariable(VarTypeString, businessKey)
	params := processStartParams{
		BusinessKey: businessKey,
		Variables:   variables,
	}

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(params)
	if err != nil {
		return fmt.Errorf("failed to send camunda message due to json error: %w", err)
	}

	url := fmt.Sprintf("process-definition/key/%s/start", c.processKey)
	_, err = c.doPostRequest(ctx, &buf, url)
	if err != nil {
		return fmt.Errorf("failed to start process for business key [%s] due to: %w", params.BusinessKey, err)
	}

	return nil
}

func (c *client) SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]CamundaVariable) error {
	buf := bytes.Buffer{}
	url := "message"
	newMessage := newMessage(messageType, businessKey, updatedVariables)
	err := json.NewEncoder(&buf).Encode(newMessage)
	if err != nil {
		return fmt.Errorf("failed to send camunda message due to json error: %w", err)
	}

	_, err = c.doPostRequest(ctx, &buf, url)
	if err != nil {
		return fmt.Errorf("failed to send message for business key [%s] due to: %w", newMessage.BusinessKey, err)
	}

	return nil
}

func (c *client) Subscribe(topicName string, handler TaskHandlerFunc, interval time.Duration) Subscription {
	sub := newSubscription(c, topicName, interval)
	sub.addHandler(handler)

	// run async fetch loop
	go sub.schedule()

	return sub
}

func (c *client) complete(ctx context.Context, taskId string, params taskCompletionParams) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(params)
	if err != nil {
		return fmt.Errorf("failed to complete camunda task due to json error: %w", err)
	}

	url := fmt.Sprintf("external-task/%s/complete", taskId)
	_, err = c.doPostRequest(ctx, &buf, url)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) fetchAndLock(param *fetchAndLock) ([]Task, error) {
	var tasks []Task
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(param)
	if err != nil {
		return tasks, fmt.Errorf("failed to fetch camunda tasks due to json error: %w", err)
	}

	url := "external-task/fetchAndLock"
	body, err := c.doPostRequest(context.Background(), &buf, url)
	if err != nil {
		return tasks, err
	}

	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return tasks, fmt.Errorf("could not unmarshal task due to: %w", err)
	}

	return tasks, nil
}

func (c *client) doPostRequest(ctx context.Context, params *bytes.Buffer, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.camundaURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, params)
	if err != nil {
		return nil, fmt.Errorf("could not create POST request due to: %w", err)
	}
	req.Header.Add(internal.HeaderContentType, "application/json")

	if c.credentials != nil {
		req.SetBasicAuth(c.credentials.User, c.credentials.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send POST request due to: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("camunda API returned Status %d with body: %v", resp.StatusCode, string(body))
	}

	return body, nil
}
