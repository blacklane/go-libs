FROM golang:1.16-alpine as baseimg

RUN apk update && \
    apk --no-cache upgrade && \
    apk --no-cache add git make && \
    apk add build-base

# First only download the dependencies, so the dependencies can be cahced before we copy the code
ARG GITHUB_TOKEN
WORKDIR /tracing/
COPY go.mod go.sum Makefile ./
RUN make gitconfig deps

WORKDIR /tracing/example/
COPY examples/go.mod examples/go.sum examples/Makefile ./
RUN make gitconfig deps

ARG BUILD_TAGS
ARG KAFKA_BOOTSTRAP_SERVERS

WORKDIR /tracing/
COPY . ./

WORKDIR /tracing/example/
COPY examples/ ./

CMD make run
