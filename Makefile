COVERAGE_OUT := $(shell test -f coverage.out && echo 1 || echo 0)
GO ?= go
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
GOFILES := $(shell find . -name "*.go")
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./...)
TESTFOLDER := $(shell find . -path "./.git" -prune -o -name "*.go" -type f -exec dirname {} +|uniq)
TESTTAGS ?= "-v"
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /examples/)

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
	@echo "Running tests with coverage report...";
	@set -eu;$(GO) mod tidy;$(GO) test -json -shuffle=on -timeout=5m -count=1 $(TESTTAGS) $(TESTFOLDER) -coverprofile=coverage.out -covermode=atomic 2>&1 | tee ./gotest-e2e.log | gotestfmt

.PHONY: coverage
coverage: ## Percentage of test coverage. If coverage <47.6%, output signal 1.
ifeq ($(COVERAGE_OUT), 0)
coverage: test
else
coverage:
endif
	@PERCENT=$$($(GO) tool cover -func=coverage.out | grep total | awk '{print $$3}'); \
	echo "coverage at: $${PERCENT}"; \
	echo $${PERCENT} | sed 's/%//' | xargs -I {} sh -c 'echo "{} < 47.6" | bc -l | grep -q 1 && exit 1 || exit 0'

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
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2; \
		if command -v goenv >/dev/null 2>&1; then \
			goenv rehash; \
		fi \
	fi
	@golangci-lint run --fix

.PHONY: misspell
misspell: ## Correct commonly misspelled English words in source code.
	misspell -w $(GOFILES)

.PHONY: misspell-check
misspell-check: ## misspell (check only).
	misspell -error $(GOFILES)

.PHONY: tools
tools: ## Install Go tools (including misspell).
	@$(GO) install github.com/client9/misspell/cmd/misspell@latest
	@$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@$(GO) install mvdan.cc/gofumpt@latest
	@if command -v goenv >/dev/null 2>&1; then \
		goenv rehash; \
	fi;

.PHONY: gotestfmt
gotestfmt: ## Install gotestfmt if not present
	@if ! command -v gotestfmt >/dev/null 2>&1; then \
		echo "Installing gotestfmt..."; \
		$(GO) install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest; \
	fi
	@if command -v goenv >/dev/null 2>&1; then \
		goenv rehash; \
	fi

.PHONY: build
build: ## Build for DEV
	@rm -f ./api
	go build -o ./api ./cmd/api

.PHONY: clean-build
clean-build: ## Build for PROD
	go mod verify
	go mod tidy
	@rm -f ./api
	CGO_ENABLED=0 go build -mod=readonly -tags netgo -trimpath -ldflags='-s -w -extldflags "-static"' -o ./api ./cmd/api
