# OpenTelemetry (OTel) example

It shows the basic usage of our `otel` module which provides [OpenTelemetry](https://opentelemetry.io/)
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

 - Call [localhost:4242](http://localhost:4242/).
 
 - Open [Jaeger UI](http://localhost:16686/search) and see the traces.

 - [Look at the code](https://github.com/blacklane/go-libs/tree/opentelemetry/otel/examples/cmd)

 - Have fun :)


## The flow

It runs an HTTP server which upon receiving a request produces an event to kafka, or fails without producing an event.
The produced event is consumed by the other service which prints the event.

Accessing [Jaeger UI](http://localhost:16686/search) you can see the flow across
the different services. 

More on OpenTracing, check the [docs](https://github.com/open-telemetry/opentelemetry-go)
