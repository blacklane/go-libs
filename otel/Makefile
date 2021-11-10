GO ?= go
GOLINTER ?= revive
GOPRIVATE ?= github.com/blacklane/*

deps: gitconfig
	GOPRIVATE=${GOPRIVATE} go mod download

.PHONY: gitconfig
gitconfig:
	git config --global url.https://${GITHUB_TOKEN}@github.com.insteadOf https://github.com

.PHONY: run
run: deps
	go run *.go

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: lint
lint:
	${GOLINTER} ./...
