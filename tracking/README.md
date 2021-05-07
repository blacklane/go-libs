# tracking

`tracking` is a working in progress lib for tracking/distributed trace resources

## Installation

```go
go get -u github.com/blacklane/go-libs/tracking
```

## Getting Started: WIP

For now there are:
 - UUID in the context
 - HTTP middleware to add/read uuid from context and request headers
 - Request depth, rout and tree path

### To see full API:

```shell script
# Install godoc
GO111MODULE=off go install golang.org/x/tools/cmd/godoc

# run godoc
godoc -http=localhost:6060
``` 

then head to:
 - [go-libs/tracking](http://localhost:6060/pkg/github.com/blacklane/go-libs/tracking/)
 - [go-libs/tracking/middleware](http://localhost:6060/pkg/github.com/blacklane/go-libs/tracking/middleware/)

# WIP: opentracing

## TODO:

 - add span to external request
 - add span to produced event headers
 - add example

1. start _jaeger_ tracer/collector to send the data to with:
````shell
docker run -d \
  --name jaeger \
  --rm -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.21
````
2. run two instances of this application

```shell
APP_NAME=first_app go run main.go
APP_NAME=second_app PORT=8001 go run main.go
```

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

5. open [Jaeger UI](http://localhost:16686/search)

6. have fun :)

## More on:
### Go
- https://github.com/opentracing/opentracing-go
- https://github.com/jaegertracing/jaeger-client-go will be replaced by Datadog Go tracer

### JS/TS/NodeJS:
- https://github.com/opentracing/opentracing-javascript
- https://github.com/jaegertracing/jaeger-client-node will be replaced by Datadog JS/TS tracer
