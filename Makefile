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

##@ Install Dependencies

.PHONY: deps
deps: ## Download go modules
	@echo "Downloading go modules..."
	go mod download

.PHONY: install
install: ## Install the uuidkey binary
	@echo "Installing uuidkey binary..."
	go install ./cmd/uuidkey

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

.PHONY: test-cli
test-cli: ## Run CLI-specific tests
	@echo "Running CLI tests..."
	go test -v ./cmd/uuidkey/...

.PHONY: coverage
coverage: ## Run tests and generate coverage report
	@echo "Running tests and generating coverage report..."
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

##@ Build & Release

.PHONY: build
build: ## Build the CLI binary
	@echo "Building uuidkey binary..."
	go build -o dist/uuidkey ./cmd/uuidkey

.PHONY: build-all
build-all: ## Build binaries for all platforms
	@echo "Building binaries for all platforms..."
	goreleaser build --snapshot --clean

.PHONY: release
release: ## Create a new release (requires version tag)
	@echo "Creating release..."
	@if [ -z "$$(git describe --tags --exact-match 2>/dev/null)" ]; then \
		echo "Error: No tag found. Please create a tag first using 'make tag VERSION=v1.2.3'"; \
		exit 1; \
	fi
	goreleaser release --clean

.PHONY: release-snapshot
release-snapshot: ## Test release process locally (doesn't publish)
	@echo "Testing release process locally..."
	goreleaser release --snapshot --clean

.PHONY: tag
tag: ## Create and push a new version tag (usage: make tag VERSION=v1.2.3)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make tag VERSION=v1.2.3"; \
		exit 1; \
	fi
	@echo "Creating tag $(VERSION)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Tag $(VERSION) created. Push with: git push origin $(VERSION)"