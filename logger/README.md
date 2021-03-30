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
    logger.WithStr("key1", "value1"),
    logger.WithStr("key2", "value2"),
    logger.WithLevel(logger.InfoLevel))

l.Info().Msg("Hello, Gophers!")
l.Debug().Msg("This Message will not appear by logger level")
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
