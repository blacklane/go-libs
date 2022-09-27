package camunda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/blacklane/go-libs/camunda/v2/internal"
)

type (
	Client interface {
		StartProcess(ctx context.Context, businessKey string, variables map[string]Variable) error
		SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]Variable) error
		Subscribe(topicName string, workerID string, handler TaskHandler, options ...func(*Subscription)) *Subscription
		DeleteTask(ctx context.Context, businessKey string) error
	}
	BasicAuthCredentials struct {
		User     string
		Password string
	}
	client struct {
		camundaURL  string
		credentials BasicAuthCredentials
		processKey  string
		httpClient  internal.HttpClient
	}
)

func NewClient(url string, processKey string, httpClient http.Client, credentials BasicAuthCredentials) Client {
	return &client{
		camundaURL:  url,
		credentials: credentials,
		processKey:  processKey,
		httpClient:  &httpClient,
	}
}

func (c *client) StartProcess(ctx context.Context, businessKey string, variables map[string]Variable) error {
	variables[businessKeyJSONKey] = NewStringVariable(VarTypeString, businessKey)
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

func (c *client) SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]Variable) error {
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

func (c *client) Subscribe(topicName string, workerID string, handler TaskHandler, options ...func(*Subscription)) *Subscription {
	sub := newSubscription(c, topicName, workerID)
	sub.addHandler(handler)

	for _, option := range options {
		option(sub)
	}

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

	req.SetBasicAuth(c.credentials.User, c.credentials.Password)

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

func (c *client) DeleteTask(ctx context.Context, businessKey string) error {
	tasks, err := c.getTasks(ctx, businessKey)
	if err != nil {
		return err
	}
	if len(tasks) == 1 {
		if err = c.deleteProcessInstance(ctx, tasks[0].ProcessInstanceId); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("found %d camunda tasks for businessKey: %s", len(tasks), businessKey)
	}
	return nil
}

func (c *client) getTasks(ctx context.Context, businessKey string) ([]Task, error) {
	url := "task"
	params := processTaskParams{
		BusinessKey: businessKey,
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, fmt.Errorf("failed to send camunda message due to json error: %w", err)
	}

	var tasks []Task
	bytes, err := c.doPostRequest(ctx, &buf, url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &tasks)
	return tasks, err
}

func (c *client) deleteProcessInstance(ctx context.Context, processInstanceId string) error {
	url := fmt.Sprintf("%s/%s/%s", c.camundaURL, "process-instance", processInstanceId)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create DELETE process-instance request due to: %w", err)
	}
	req.SetBasicAuth(c.credentials.User, c.credentials.Password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send DELETE process-instance request due to: %w", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("camunda delete process-instance API returned Status %d", resp.StatusCode)
	}
	return nil
}
