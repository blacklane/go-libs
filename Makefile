GO ?= go
TIMEOUT ?= 60
BUILD_TAGS ?=

TOOLS_MOD_DIR := ./internal/tools
ALL_DOCS := $(shell find . -name '*.md' -type f | sort)
ALL_GO_MOD_DIRS := $(shell find . -name go.mod | xargs dirname | sort)
ALL_COVERAGE_MOD_DIRS := $(shell find . -type f -name go.mod | xargs dirname | egrep -v '^./example|^$(TOOLS_MOD_DIR)' | sort)
GO_MOD_DIRS := $(filter-out $(TOOLS_MOD_DIR), $(ALL_GO_MOD_DIRS))

# Tools

.PHONY: tools

TOOLS = $(CURDIR)/.tools

$(TOOLS):
	@mkdir -p $@
$(TOOLS)/%: | $(TOOLS)
	cd $(TOOLS_MOD_DIR) && \
	$(GO) build -o $@ $(PACKAGE)

GOLANGCI_LINT = $(TOOLS)/golangci-lint
$(GOLANGCI_LINT): PACKAGE=github.com/golangci/golangci-lint/cmd/golangci-lint

MISSPELL = $(TOOLS)/misspell
$(MISSPELL): PACKAGE=github.com/client9/misspell/cmd/misspell

MULTIMOD = $(TOOLS)/multimod
$(MULTIMOD): PACKAGE=go.opentelemetry.io/build-tools/multimod

STRINGER = $(TOOLS)/stringer
$(STRINGER): PACKAGE=golang.org/x/tools/cmd/stringer

GOCOVMERGE = $(TOOLS)/gocovmerge
$(GOCOVMERGE): PACKAGE=github.com/wadey/gocovmerge

tools: $(GOLANGCI_LINT) $(MISSPELL) $(GOCOVMERGE) $(STRINGER) $(MULTIMOD)

# Linting

.PHONY: golangci-lint golangci-lint-fix
golangci-lint-fix: ARGS=--fix
golangci-lint-fix: golangci-lint
golangci-lint: $(GO_MOD_DIRS:%=golangci-lint/%)
golangci-lint/%: DIR=$*
golangci-lint/%: | $(GOLANGCI_LINT)
	@echo 'golangci-lint $(if $(ARGS),$(ARGS) ,)$(DIR)' \
		&& cd $(DIR) \
		&& $(GOLANGCI_LINT) run --allow-serial-runners $(ARGS)

.PHONY: go-mod-tidy
go-mod-tidy: $(ALL_GO_MOD_DIRS:%=go-mod-tidy/%)
go-mod-tidy/%: DIR=$*
go-mod-tidy/%:
	@echo "$(GO) mod tidy in $(DIR)" \
		&& cd $(DIR) \
		&& $(GO) mod tidy

.PHONY: misspell
misspell: | $(MISSPELL)
	@$(MISSPELL) -w $(ALL_DOCS)

.PHONY: lint
lint: go-mod-tidy golangci-lint misspell

# Build

.PHONY: generate build

generate: $(GO_MOD_DIRS:%=generate/%)
generate/%: DIR=$*
generate/%: | $(STRINGER)
	@echo "$(GO) generate $(DIR)/..." \
		&& cd $(DIR) \
		&& PATH="$(TOOLS):$${PATH}" $(GO) generate ./...

build: generate $(GO_MOD_DIRS:%=build/%) $(GO_MOD_DIRS:%=build-tests/%)
build/%: DIR=$*
build/%:
	@echo "$(GO) build ${BUILD_TAGS} $(DIR)/..." \
		&& cd $(DIR) \
		&& $(GO) build ${BUILD_TAGS} ./...

build-tests/%: DIR=$*
build-tests/%:
	@echo "$(GO) build tests ${BUILD_TAGS} $(DIR)/..." \
		&& cd $(DIR) \
		&& $(GO) list ./... \
		| grep -v third_party \
		| xargs $(GO) test ${BUILD_TAGS} -vet=off -run xxxxxMatchNothingxxxxx >/dev/null

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test
test-default test-race: ARGS=-race
test-bench:   ARGS=-run=xxxxxMatchNothingxxxxx -test.benchtime=1ms -bench=.
test-short:   ARGS=-short
test-verbose: ARGS=-v
$(TEST_TARGETS): test
test: $(GO_MOD_DIRS:%=test/%)
test/%: DIR=$*
test/%:
	@echo "$(GO) test ${BUILD_TAGS} -timeout $(TIMEOUT)s $(ARGS) $(DIR)/..." \
		&& cd $(DIR) \
		&& $(GO) test ${BUILD_TAGS} -timeout $(TIMEOUT)s $(ARGS) ./...

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = coverage.out
.PHONY: test-coverage
test-coverage: $(ALL_COVERAGE_MOD_DIRS:%=test-coverage/%) | $(GOCOVMERGE)
	@printf "" > coverage.txt \
		&& $(GOCOVMERGE) $$(find . -name $(COVERAGE_PROFILE)) > coverage.txt

test-coverage/%: DIR=$*
test-coverage/%:
	@set -e; \
		CMD="$(GO) test ${BUILD_TAGS} -race -covermode=$(COVERAGE_MODE) -coverprofile=$(COVERAGE_PROFILE)"; \
		echo "$(DIR)" | grep -q 'test$$' \
		&& CMD="$$CMD -coverpkg=go.opentelemetry.io/contrib/$$( dirname "$(DIR)" | sed -e "s/^\.\///g" )/..."; \
		echo "$$CMD $(DIR)/..."; \
		cd "$(DIR)" \
		&& $$CMD ./... \
		&& $(GO) tool cover -html=coverage.out -o coverage.html;

# Workspace

go-work-init:
	@rm -f go.work*
	@$(GO) work init
	@for mod in $(ALL_GO_MOD_DIRS); do \
		$(GO) work use $$mod; \
	done

# Releasing

.PHONY: prerelease
prerelease: | $(MULTIMOD)
	@[ "${MODSET}" ] || ( echo ">> env var MODSET is not set"; exit 1 )
	$(MULTIMOD) verify && $(MULTIMOD) prerelease -m ${MODSET}

COMMIT ?= "HEAD"
.PHONY: add-tags
add-tags: | $(MULTIMOD)
	@[ "${MODSET}" ] || ( echo ">> env var MODSET is not set"; exit 1 )
	$(MULTIMOD) verify && $(MULTIMOD) tag -m ${MODSET} -c ${COMMIT}

VERSION ?=
TAGS ?= $(shell git tag -l | grep "$(VERSION)")
push-tags:
	@[ "${VERSION}" ] || ( echo ">> env var VERSION is not set"; exit 1 )
	for t in $(TAGS); do git push origin "$$t"; done