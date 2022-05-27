FROM golang:1.18-alpine

RUN apk update && \
    apk --no-cache upgrade && \
    apk --no-cache add \
    build-base \
    git \
    make

WORKDIR /src
COPY . ./

ENTRYPOINT [ "sh", "-c" ]