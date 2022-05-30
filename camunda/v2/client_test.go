package camunda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var basicAuth = BasicAuthCredentials{User: "Bernd", Password: "BerndsPassword"}

const testWorkerID = "test-worker"

func TestClient_StartProcess(t *testing.T) {
	client, mockHttpClient := newTestClient()
	businessKey := uuid.New().String()
	variables := map[string]Variable{
		"test-var": {Type: VarTypeInteger, Value: 42},
	}

	body := ioutil.NopCloser(bytes.NewReader([]byte("{}")))
	mockHttpClient.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		request := args.Get(0).(*http.Request)
		requestBody := parseStartProcessRequestBody(t, request.Body)

		assert.Equal(t, businessKey, requestBody.BusinessKey)
		assert.Equal(t, len(variables), len(requestBody.Variables))
		assert.Equal(t, VarTypeInteger, requestBody.Variables["test-var"].Type)
		assert.Equal(t, 42, int(requestBody.Variables["test-var"].Value.(float64)))
	}).Return(&http.Response{
		StatusCode: 200,
		Body:       body,
	})

	// act
	err := client.StartProcess(context.Background(), businessKey, variables)

	// assert
	mockHttpClient.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestClient_StartProcess_UsesBasicAuth(t *testing.T) {
	client, mockHttpClient := newTestClient()
	businessKey := uuid.New().String()
	variables := map[string]Variable{}

	url := fmt.Sprintf("/process-definition/key/%s/start", client.processKey)
	body := ioutil.NopCloser(bytes.NewReader([]byte("{}")))
	mockHttpClient.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		request := args.Get(0).(*http.Request)
		assert.Equal(t, url, request.URL.Path)
		assert.Equal(t, client.camundaURL, fmt.Sprintf("%s://%s", request.URL.Scheme, request.URL.Host))
		user, password, ok := request.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, basicAuth.User, user)
		assert.Equal(t, basicAuth.Password, password)
	}).Return(&http.Response{
		StatusCode: 200,
		Body:       body,
	})

	// act
	err := client.StartProcess(context.Background(), businessKey, variables)

	// assert
	mockHttpClient.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestClient_SendMessage(t *testing.T) {
	client, mockHttpClient := newTestClient()

	messageType := "some-message"
	businessKey := uuid.New().String()
	variables := map[string]Variable{
		"test-var": {Type: VarTypeString, Value: "Zweiundvierzig"},
	}

	body := ioutil.NopCloser(bytes.NewReader([]byte("{}")))
	mockHttpClient.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		request := args.Get(0).(*http.Request)
		assert.Equal(t, "/message", request.URL.Path)

		requestBody := &message{}
		byteBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to parse reqeuest body due to %s", err)
		}
		err = json.Unmarshal(byteBody, requestBody)
		if err != nil {
			t.Fatalf("failed to parse reqeuest body due to %s", err)
		}
		assert.Equal(t, businessKey, requestBody.BusinessKey)
		assert.Equal(t, messageType, requestBody.MessageName)
		assert.Equal(t, "Zweiundvierzig", requestBody.ProcessVariables["test-var"].Value)
	}).Return(&http.Response{
		StatusCode: 200,
		Body:       body,
	})

	// act
	err := client.SendMessage(context.Background(), messageType, businessKey, variables)

	// assert
	mockHttpClient.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestCamundaClient_SendMessage_HandlesRequestError(t *testing.T) {
	client, mockHttpClient := newTestClient()

	messageType := "some-message"
	businessKey := uuid.New().String()
	variables := map[string]Variable{}

	body := ioutil.NopCloser(bytes.NewReader([]byte("{}")))
	mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 500,
		Body:       body,
	})

	// act
	err := client.SendMessage(context.Background(), messageType, businessKey, variables)

	// assert
	mockHttpClient.AssertExpectations(t)
	assert.Error(t, err)
}

func TestSubscription_Complete(t *testing.T) {
	client, mockHttpClient := newTestClient()

	ctx := context.Background()
	taskID := uuid.New().String()
	params := taskCompletionParams{
		WorkerID: testWorkerID,
	}

	body := ioutil.NopCloser(bytes.NewReader([]byte("{}")))
	mockHttpClient.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		request := args.Get(0).(*http.Request)
		assert.Equal(t, fmt.Sprintf("/external-task/%s/complete", taskID), request.URL.Path)

		requestBody := &taskCompletionParams{}
		byteBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to parse reqeuest body due to %s", err)
		}
		err = json.Unmarshal(byteBody, requestBody)
		if err != nil {
			t.Fatalf("failed to parse reqeuest body due to %s", err)
		}
		assert.Equal(t, params.WorkerID, requestBody.WorkerID)
	}).Return(&http.Response{
		StatusCode: 200,
		Body:       body,
	})

	// act
	err := client.complete(ctx, taskID, params)

	assert.Nil(t, err)
}

func TestSubscription_FetchAndLock(t *testing.T) {
	client, mockHttpClient := newTestClient()

	taskID := uuid.New().String()
	topic := "test-topic"
	params := &fetchAndLock{
		WorkerID: testWorkerID,
		MaxTasks: defaultMaxTasksFetch,
	}

	bodyString := fmt.Sprintf(`[{"id":"%s", "topicName": "%s"}]`, taskID, topic)
	body := ioutil.NopCloser(bytes.NewReader([]byte(bodyString)))
	mockHttpClient.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		request := args.Get(0).(*http.Request)
		assert.Equal(t, "/external-task/fetchAndLock", request.URL.Path)

		requestBody := &fetchAndLock{}
		byteBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to parse reqeuest body due to %s", err)
		}
		err = json.Unmarshal(byteBody, requestBody)
		if err != nil {
			t.Fatalf("failed to parse reqeuest body due to %s", err)
		}
		assert.Equal(t, params.WorkerID, requestBody.WorkerID)
		assert.Equal(t, params.MaxTasks, requestBody.MaxTasks)
	}).Return(&http.Response{
		StatusCode: 200,
		Body:       body,
	})

	// act
	tasks, err := client.fetchAndLock(params)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, taskID, tasks[0].ID)
	assert.Equal(t, topic, tasks[0].TopicName)
}

func newTestClient() (*client, *MockHttpClient) {
	mockHttpClient := &MockHttpClient{}
	client := &client{
		camundaURL:  "http://testing.local",
		processKey:  "TestingRides",
		httpClient:  mockHttpClient,
		credentials: basicAuth,
	}
	return client, mockHttpClient
}

func parseStartProcessRequestBody(t *testing.T, body io.ReadCloser) processStartParams {
	requestBody := &processStartParams{}
	byteBody, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatalf("failed to parse reqeuest body due to %s", err)
	}
	err = json.Unmarshal(byteBody, requestBody)
	if err != nil {
		t.Fatalf("failed to parse reqeuest body due to %s", err)
	}

	return *requestBody
}

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(r *http.Request) (*http.Response, error) {
	args := m.Called(r)
	return args.Get(0).(*http.Response), nil
}
