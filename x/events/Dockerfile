FROM golang:1.14.3-alpine as baseimg

ARG BUILD_TAGS="-tags musl,integration"

RUN apk update && \
    apk --no-cache upgrade && \
    apk --no-cache add git make && \
    apk add build-base

# First only download the dependencies, so the dependencies can be cahced before we copy the code
COPY ./go.mod ./go.sum ./Makefile /app/
WORKDIR /app/

RUN make deps install_golint
COPY . ./

CMD make check test
