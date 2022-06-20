# Middleware

TODO: write the readme

## Installation

```go
go get -u github.com/blacklane/go-libs/middleware
```

## Retry: Getting Started

`retry` is a middleware which for x/events. It will retry with exponential back-off algorithm until it exceeds max retry, which is passed by argument to `middleware.Retry`.

```go
maxRetryCount := 10

hb := events.HandlerBuilder{}
hb.AddHandler(handler)
hb.UseMiddleware(middleware.Retry(maxRetryCount))
```

If the handler returns a RetriableError, then the middleware will retry.

```go
type Handler struct {}

func (h Handler) Handle(ctx context.Context, e events.Event) error {
	err := processEvent(e)
	if err != nil {
		return RetriableError{
			Retriable: true,
			Err:       fmt.Errorf("fail to process event: %w", err),
		}
	}

	return nil
}
```

## HTTPWithBodyFilter Middleware

`HTTPWithBodyFilter` It will log body from the request, this middleware also accepts a string slice with Keys to 
filter sensitive data.

```go
filterKeys := []string{"phone", "account_number", "sensitive_data"}

hb := rest.HandlerBuilder{}
hb.AddHandler(handler)
hb.UseMiddleware(middleware.HTTPWithBodyFilter("serviceName", "handlerName", "/path", filterKeys, Logger))
```

Example:

Input
```json
{
  "name": "John",
  "phone": "1234524",
  "account_number": "DE21343221233124"
}
```
Output in logs
```json
{
  "name": "John",
  "phone": "[FILTERED]",
  "account_number": "[FILTERED]"
}
```

Notes:
* You can pass ```nil``` as filterKeys to log all body data, this is going to use ```internal.DefaultKeys``` for filtering.
* You can pass a ```non-existing-key``` as filterKeys to log all body data without applying ```internal.DefaultKeys```.

Tips:
* You can create a *filterKeys.yaml* file to define keys to filter by path, and then apply the middleware
on each endpoint with the required data.

Example:

filterKeys.yaml: 
```yaml
"/path1": ["cardholder_name","email", "new_passenger_phone", "phone", "first_name"],
"/path2": ["name","address"]

```

You can load the data in the config:
```go
type Config struct {
AppName    string `env:"APP_NAME" envDefault:"name"`
FilterKeys map[string][]string
...
}

func (cfg Config) GetFilterKeys() map[string][]string {
//get filters
yfile, err := ioutil.ReadFile("filtersKeys.yaml")
if err != nil {
    log.Panic(err)
}
filterData := make(map[string][]string)
err = yaml.Unmarshal(yfile, &filterData)
if err != nil {
    log.Panic(err)
}

return filterData
}
```

And apply it by endpoint:

```go
package example

import (
	logmiddleware "github.com/blacklane/go-libs/middleware"
)
logmiddleware.HTTPWithBodyFilter("app-name", "handler-name", "/path1", cfg.FilterKeys["/path1"], cfg.Logger),
```



## Camunda Subscription Handler: Getting Started

The Camunda subscription handler is a middleware which allows the injection of the `context` from the service in ot the Camunda subscription handler.

```go
func CamundaSubscriptionsAddLogger(log logger.Logger) func(camunda.TaskHandler) camunda.TaskHandler {
	return func(next camunda.TaskHandler) camunda.TaskHandler {
		return camunda.TaskHandlerFunc(func(ctx context.Context, completeFunc camunda.TaskCompleteFunc, t camunda.Task) {
			log = log.With().Logger()
			context := log.WithContext(ctx)
			next.Handle(context, completeFunc, t)
		})
	}
}
```
