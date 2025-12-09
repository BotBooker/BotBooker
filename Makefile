GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /examples/)
GOFILES := $(shell find . -name "*.go")
TESTFOLDER := $(shell find . -name "*_test.go" -type f -exec dirname {} +)
TESTTAGS ?= ""

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "\033[1;3;34mBotBooker Go.\033[0m\n"
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_0-9\/-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run tests to verify code functionality.
test: gotestfmt
	@echo "Running tests with coverage report..."
	@set -euo pipefail
	@$(GO) test -json -v -shuffle=on -timeout=5m -count=1 -tags cmd ./cmd/... -coverprofile=coverage.out -covermode=atomic 2>&1 | tee ./gotest-e2e.log | gotestfmt

.PHONY: fmt
fmt: ## Ensure consistent code formatting.
	@$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check: ## format (check only).
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: vet
vet: ## Examine packages and report suspicious constructs if any.
	@$(GO) vet $(VETPACKAGES)

.PHONY: lint
lint: ## Inspect source code for stylistic errors or potential bugs.
	@golangci-lint run --fix

.PHONY: misspell
misspell: ## Correct commonly misspelled English words in source code.
	misspell -w $(GOFILES)

.PHONY: misspell-check
misspell-check: ## misspell (check only).
	misspell -error $(GOFILES)

.PHONY: tools
tools: ## Install tools (golint and misspell).
	@if [ $(GO_VERSION) -gt 15 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
		$(GO) install github.com/client9/misspell/cmd/misspell@latest; \
	elif [ $(GO_VERSION) -lt 16 ]; then \
		$(GO) install golang.org/x/lint/golint; \
		$(GO) install github.com/client9/misspell/cmd/misspell; \
	fi
	@if command -v goenv >/dev/null 2>&1; then \
		goenv rehash; \
	fi

.PHONY: gotestfmt
gotestfmt: ## Install gotestfmt if not present
	@if ! command -v gotestfmt >/dev/null 2>&1; then \
		echo "Installing gotestfmt..."; \
		$(GO) install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest; \
	fi
	@if command -v goenv >/dev/null 2>&1; then \
		goenv rehash; \
	fi
