package camunda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/blacklane/go-libs/logger"

	"github.com/blacklane/go-libs/x/camunda/internal"
)

type (
	Client interface {
		StartProcess(ctx context.Context, businessKey string, variables map[string]CamundaVariable) error
		SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]CamundaVariable)
		//Subscribe(topicName string, handler TaskHandlerFunc, interval time.Duration) Subscription
	}
	camundaClient struct {
		camundaURL  string
		credentials BasicAuthCredentials
		processKey  string
		httpClient  HttpClient
		log         logger.Logger
	}
)

func NewClient(log logger.Logger, url string, processKey string, credentials BasicAuthCredentials) Client {
	return &camundaClient{
		camundaURL:  url,
		credentials: credentials,
		processKey:  processKey,
		httpClient:  &http.Client{},
		log:         log,
	}
}

func (c *camundaClient) StartProcess(ctx context.Context, businessKey string, variables map[string]CamundaVariable) error {
	params := processStartParams{
		BusinessKey: businessKey,
		Variables:   variables,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(params)
	if err != nil {
		c.log.Err(err).Msg("failed to send camunda message due to json error")
		return err
	}

	url := fmt.Sprintf("process-definition/key/%s/start", c.processKey)
	_, err = c.doPostRequest(ctx, buf, url)
	if err != nil {
		c.log.Err(err).
			Msgf("failed to start process for business key - %s", params.BusinessKey)
		return err
	}
	c.log.Info().Msgf("process has been started with business key '%s'", params.BusinessKey)
	return nil
}

func (c *camundaClient) SendMessage(ctx context.Context, messageType string, businessKey string, updatedVariables map[string]CamundaVariable) {
	buf := new(bytes.Buffer)
	url := "message"
	newMessage := newMessage(messageType, businessKey, updatedVariables)
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

func (c *camundaClient) doPostRequest(ctx context.Context, params *bytes.Buffer, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.camundaURL, endpoint)

	req, err := http.NewRequest("POST", url, params)
	if err != nil {
		c.log.Err(err).Msgf("Could not create POST request")
		return nil, err
	}
	req.Header.Add(internal.HeaderContentType, "application/json")
	req.SetBasicAuth(c.credentials.User, c.credentials.Password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Err(err).Msgf("Could not send POST request")
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		c.log.Error().Str(internal.LogFieldURL, url).Msgf("Status: %v, Body: %v", resp.StatusCode, string(body))
		return nil, errors.New(fmt.Sprintf("Camunda API returned Status %v", resp.StatusCode))
	}

	return body, nil
}
