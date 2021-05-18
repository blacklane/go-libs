FROM golang:1.16-alpine as baseimg

ARG BUILD_TAGS="-tags musl,integration"

RUN apk update && \
    apk --no-cache upgrade && \
    apk --no-cache add git make && \
    apk add build-base

# First only download the dependencies, so the dependencies can be cahced before we copy the code
WORKDIR /app/
COPY go.mod ./go.sum ./Makefile ./

ARG GITHUB_TOKEN

RUN make gitconfig deps
COPY . ./

CMD go run *.go
