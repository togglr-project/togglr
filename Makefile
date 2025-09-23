NAMESPACE=togglr
GOCI_VERSION=v1.64.8
MOCKERY_VERSION=v3.5.4
TOOLS_DIR=dev/tools
TOOLS_DIR_ABS=${PWD}/${TOOLS_DIR}
BIN_OUTPUT_DIR=bin
GOLANGCI_LINT=${TOOLS_DIR_ABS}/golangci-lint
MOCKERY=${TOOLS_DIR_ABS}/mockery
GOCMD=go
GOBUILD=$(GOCMD) build
GOPROXY=https://proxy.golang.org,direct
TOOL_VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
TOOL_BUILD_TIME=$(shell date '+%Y-%m-%dT%H:%M:%SZ%Z')
OS=$(shell uname -s)

LD_FLAGS="-w -s -X 'github.com/togglr-project/togglr/internal/version.Version=${TOOL_VERSION}' -X 'github.com/togglr-project/togglr/internal/version.BuildTime=${TOOL_BUILD_TIME}'"

RED="\033[0;31m"
GREEN="\033[1;32m"
YELLOW="\033[0;33m"
NOCOLOR="\033[0m"

.DEFAULT_GOAL := help

#
# Extra targets
#
-include dev/dev.mk

#
# Local targets
#

.PHONY: help
help: ## Print this message
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: .install-linter
.install-linter:
	@[ -f $(GOLANGCI_LINT) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR) $(GOCI_VERSION)

.PHONY: .install-mockery
.install-mockery:
	@[ -f $(MOCKERY) ] || GOBIN=$(TOOLS_DIR_ABS) go install github.com/vektra/mockery/v3@$(MOCKERY_VERSION)

.PHONY: setup
setup: .install-linter .install-mockery ## Setup development environment
	@echo "\nCreate .env files in dev/ directory"
	@cp dev/config.env.example dev/config.env
	@cp dev/compose.env.example dev/compose.env
	@cp dev/platform.env.example dev/platform.env

	@echo
	@if [ $$? -ne 0 ] ; then \
		@echo -e ${RED}"FAIL"${NOCOLOR} ; \
		exit 1 ; \
	fi
	@echo ${GREEN}"OK"${NOCOLOR}

.PHONY: lint
lint: .install-linter ## Run linter
	@$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: test
test: ## Run unit tests
	@go test -cover -coverprofile=coverage.out -v ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: test-integration
test-integration: ## Run unit and integration tests
	@go test -tags integration -cover -coverprofile=coverage.out -v ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: build
build: ## Build binary
	@echo "\nBuilding binary..."
	@echo
	go env -w GOPROXY=${GOPROXY}
	go env -w GOPRIVATE=${GOPRIVATE}

	CGO_ENABLED=0 $(GOBUILD) -trimpath -tags timetzdata -ldflags=$(LD_FLAGS) -o ${BIN_OUTPUT_DIR}/app .

.PHONY: mocks
mocks: .install-mockery ## Generate mocks with mockery
	@rm -rf test_mocks
	$(MOCKERY)
	@mv test_mocks/internal test_mocks/internal_

.PHONY: format
format: ## Format go code
	go fmt ./...

.PHONY: generate-backend
generate-backend: ## Generate backend by OpenAPI specification
	@docker run --rm \
      --volume ".:/workspace" \
      ghcr.io/ogen-go/ogen:latest --target workspace/internal/generated/server --clean workspace/specs/server.yml

.PHONY: generate-sdk-backend
generate-sdk-backend: ## Generate SDK backend by OpenAPI specification
	@docker run --rm \
      --volume ".:/workspace" \
      ghcr.io/ogen-go/ogen:latest --target workspace/internal/generated/sdkserver --clean workspace/specs/sdk.yml
