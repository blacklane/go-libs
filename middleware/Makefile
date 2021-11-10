GO ?= go
GOLINTER ?= revive
GOPRIVATE ?= github.com/blacklane/*

# Use a github token or similar to access private repos
# In this example it's got no effect rather than changing from ssh to https
.PHONY: gitconfig
gitconfig:
	@git config --global url.https://${GITHUB_TOKEN}@github.com.insteadOf https://github.com

.PHONY: deps
deps:
	GOPRIVATE=${GOPRIVATE} ${GO} mod download

.PHONY: test
test: deps
	${GO} test -race -cover ./...

.PHONY: lint
lint:
	${GOLINTER} ./...
