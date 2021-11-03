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

.PHONY: goimports
goimports:
	ifndef file
		$(error file is not set. Usage: make goimports file=path/to/file)
	endif
	echo @(file)

.PHONY: lint
lint:
	golint ./...
	staticcheck ./...

.PHONY: install_golint
install_linters:
	go get -u golang.org/x/lint/golint
	go get -u honnef.co/go/tools/cmd/staticcheck@latest
