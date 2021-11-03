# OpenTelemetry (OTel) tracing example

It shows the basic usage of our otel package which provides [OpenTelemetry](https://opentelemetry.io/)
middleware for [HTTP](https://golang.org/pkg/net/http/#Handler)
and [Events](https://github.com/blacklane/go-libs/blob/master/x/events/events.go#L14) 
handlers as well as the necessary setup to use OTel.

## Running

- Run the dependencies

```shell
make compose-dependencies
```

- Run the Go examples

```shell
make run-http
```

```shell
make run-events
```

 - You call [localhost:4242](http://localhost:4242/) directly, or use our graphql (see below).

 - Graphql using branch `opentelemetry`:

```shell
# Run a redis on docker so graphql won't be logging thousands of errors :/
docker run --rm --name graphql-redis -d -p 6379:6379 redis:alpine
```

```shell
source .env.example
KAFKA_USE_SASL=false REDIS_HOST=localhost npm run start:dev
```

 - Execute the _tracing_ query a few times.
```graphql
query {
  tracing
}
```

 - Open [Jaeger UI](http://localhost:16686/search) and see the traces.

 - [Look at the code](https://github.com/blacklane/go-libs/tree/opentelemetry/otel/examples/cmd)

 - Have fun :)


## The flow

It simulates an HTTP server which upon receiving a request produces an event to kafka, or fails without producing an event.
The produced event is consumed by another service which prints the event.

Accessing [Jaeger UI](http://localhost:16686/search) you can see the flow across
the different services. There is also an automatically generated [service map](http://localhost:16686/dependencies).

Optionally the `teacing` graphql query can call the HTTP server.

## More on OpenTracing:
### Go:
- https://github.com/open-telemetry/opentelemetry-go

### JS/TS/NodeJS:
- https://github.com/open-telemetry/opentelemetry-js
