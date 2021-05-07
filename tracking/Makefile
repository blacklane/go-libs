.PHONY: download lint

GO ?= go
GOLINT ?= golint
STATICCHECK ?= staticcheck

# Use a github token to access private repos
gitconfig:
	@git config --global url.https://${GITHUB_TOKEN}@github.com.insteadOf https://github.com

deps:
	${GO} mod download

test: deps
	${GO} test -v -race -cover ./...

lint:
	${GOLINT} ./...
	${STATICCHECK} ./...
