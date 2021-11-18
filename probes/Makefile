GO ?= go
GOLINTER ?= revive

.PHONY: deps
deps:
	GOPRIVATE=${GOPRIVATE} ${GO} mod download

.PHONY: test
test: deps
	${GO} test -race -cover ./...

.PHONY: lint
lint:
	${GOLINTER} ./...
