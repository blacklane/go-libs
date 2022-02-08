# Camunda API client

This library is a wrapper for the [camunda API](https://docs.camunda.org/manual/7.15/reference/rest/).
It can connect to any camunda engine deployment and supports a subset of operations from the official camunda API.

## Supported Operations
* starting a process instance
* sending a message (and updating variables)
* subscribing to an external task queue
* using basic authentication

## Installation
```shell
go get -u github.com/blacklane/go-libs/camunda/v2
```

## Getting Started

Create a new API client:
```go
url         := "http://localhost:8080"
processKey  := "example-process"
credentials := camunda.BasicAuthCredentials{...}
client := camunda.NewClient(url, processKey, http.Client{}, credentials)
```

Starting a new process instance:
```go
businessKey := uuid.New().String()
variables := map[string]camunda.Variable{}
err := client.StartProcess(context.Background(), businessKey, variables)
if err != nil {
	log.Err(err).Msg("Failed to start process")
}
```

For complete examples check out the [/examples](/examples) directory.
