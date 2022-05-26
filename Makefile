GO ?= go
BUILD_TAGS ?=

work-init:
	@rm -f go.work*
	@$(GO) work init
	@for mod in $(shell find . -name go.mod | xargs dirname); do \
		$(GO) work use $$mod; \
	done
	
test-all:
	@$(GO) test ${BUILD_TAGS} -v $(shell go list -f '{{.Dir}}/...' -m | xargs)
