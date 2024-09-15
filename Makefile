MAKEFLAGS += --no-print-directory

.PHONY: all
all: generate

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display the list of targets and their descriptions
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1mUsage:\033[0m\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } \
		/^###/ { printf "  \033[90m%s\033[0m\n", substr($$0, 4) }' $(MAKEFILE_LIST)

##@ Tooling 

.PHONY: install-devbox
install-devbox: ## Install Devbox
	@echo "Installing Devbox..."
	@curl -fsSL https://get.jetify.dev | bash

.PHONY: devbox-update
devbox-update: ## Update Devbox
	@devbox update

.PHONY: devbox
devbox: ## Run Devbox shell
	@devbox shell

##@ Installation

.PHONY: install
install: ## Download go modules
	@echo "Downloading go modules..."
	go mod download

##@ Development

.PHONY: fmt
fmt: ## Run go fmt
	@echo "Running go fmt..."
	go fmt ./...

.PHONY: generate
generate: ## Generate and embed go documentation into README.md
	@echo "Generating and embedding go documentation into README.md..."
	go generate ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run ./...

##@ Benchmarking, Testing, & Coverage

.PHONY: bench
bench: ## Run Go benchmarks
	@echo "Running go benchmarks..."
	go test ./... -tags=bench -bench=.

.PHONY: test
test: ## Run Go tests
	@echo "Running go tests..."
	go test ./... -tags=test

.PHONY: coverage
coverage: ## Run tests and generate coverage report
	@echo "Running tests and generating coverage report..."
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...