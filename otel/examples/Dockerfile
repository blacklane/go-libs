FROM golang:1.17-alpine as baseimg

RUN apk update && \
    apk --no-cache upgrade && \
    apk --no-cache add git make && \
    apk add build-base

ARG GITHUB_TOKEN
# First only download the dependencies, so the dependencies can be cahced before we copy the code
WORKDIR /otel/
COPY go.mod go.sum Makefile ./
RUN make gitconfig deps
COPY . ./

WORKDIR /otel/example/
COPY examples/go.mod examples/go.sum examples/Makefile ./
RUN make gitconfig deps

ARG BUILD_TAGS
ARG KAFKA_BOOTSTRAP_SERVERS

WORKDIR /otel/
COPY . ./

WORKDIR /otel/example/
COPY examples/ ./

CMD make run
