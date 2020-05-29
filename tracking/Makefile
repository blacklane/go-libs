.PHONY: download lint

GO ?= go
GOLINT ?= golint
STATICCHECK ?= staticcheck

download:
	${GO} mod download

test:
	${GO} test -race ./...

lint:
	${GOLINT} ./...
	${STATICCHECK} ./...
