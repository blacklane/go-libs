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
