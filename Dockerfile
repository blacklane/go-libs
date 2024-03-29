FROM golang:1.17-alpine3.13

RUN apk update && \
    apk --no-cache upgrade && \
    apk --no-cache add \
    build-base \
    git \
    make

ARG GITHUB_TOKEN
RUN git config --global url.https://${GITHUB_TOKEN}:@github.com.insteadOf https://github.com

WORKDIR /src
COPY . ./
