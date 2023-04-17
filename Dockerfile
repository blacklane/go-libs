FROM golang:1.17-alpine3.13

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot \
    && apk update && \
    apk --no-cache upgrade && \
    apk --no-cache add \
    build-base \
    git \
    make

USER nonroot

RUN git config --global credential.helper 'store --file /run/secrets/github_token'

WORKDIR /src
COPY . ./
