# Tracing example

It shows the basic usage of our tracing package which provides [OpenTracing](https://opentracing.io/)
middleware for [HTTP](https://golang.org/pkg/net/http/#Handler)
and [Events](https://github.com/blacklane/go-libs/blob/master/x/events/events.go#L14) 
handlers as well as functions to get the opntracing.Span from the context.

1. Start Jaeger Tracer and Kafka using docker compose:
```shell
make compose-up
```

1. Run the example
```shell
make run
```

1. Start a flow by accessing [localhost:4242](http://localhost:4242/)

1. open [Jaeger UI](http://localhost:16686/search)

1. [look at the code](https://github.com/blacklane/go-libs/blob/all/opentracing/tracking/examples/main.go#L59)

1. have fun :)

## The flow

It simulates an HTTP server which upon receiving a request produces an event to kafka, or fails without producing an event.
The produced event is consumed by another service which prints the event.

Accessing [Jaeger UI](http://localhost:16686/search) you can see the flow across
the different services. There is also an automatically generated [service map](http://localhost:16686/dependencies).

## TODO: integrate graphql
3. run graphql using branch `opentracing`
   Run a redis on docker so graphql won't be logging thousands of errors :/
```shell
docker run --rm --name graphql-redis -d -p 6379:6379 redis:alpine
npm run start:dev
```

4. execute the _tracing_ query a few times
```graphql
query {
  tracing
}
```


## More on:
### Go:
- https://github.com/opentracing/opentracing-go
- https://github.com/jaegertracing/jaeger-client-go will be replaced by Datadog Go tracer

### JS/TS/NodeJS:
- https://github.com/opentracing/opentracing-javascript
- https://github.com/jaegertracing/jaeger-client-node will be replaced by Datadog JS/TS tracer
