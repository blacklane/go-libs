# logger

`logger` is a working in progress logger for Blacklane. It's a wrapper around 
[zerolog](https://github.com/rs/zerolog)

## Installation

```go
go get -u github.com/blacklane/go-libs/logger
```

## Getting Started: WIP

```go

l := logger.New(
    os.Stdout,
    "example",
    Str("key1", "value1"),
    Str("key2", "value2"))

l.Info().Msg("Hello, Gophers!")
```

### To see full API:

```shell script
# Install godoc
GO111MODULE=off go install golang.org/x/tools/cmd/godoc

# run godoc
godoc -http=localhost:6060
``` 

then head to:
 - [go-libs/logger](http://localhost:6060/pkg/github.com/blacklane/go-libs/logger/)
 - [go-libs/logger/middleware](http://localhost:6060/pkg/github.com/blacklane/go-libs/logger/middleware/)
