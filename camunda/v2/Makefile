.PHONY: download lint

GO ?= go
GOLINT ?= golint
GOPRIVATE ?= github.com/blacklane/*
STATICCHECK ?= staticcheck

# Use a github token or similar to access private repos
# In this example it's got no effect rather than changing from ssh to https
gitconfig:
	@git config --global url.https://${GITHUB_TOKEN}@github.com.insteadOf https://github.com

deps:
	GOPRIVATE=${GOPRIVATE} ${GO} mod download

test: deps
	${GO} test -race -cover ./...

lint:
	${GOLINT} ./...
	${STATICCHECK} ./...
