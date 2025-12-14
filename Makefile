.PHONY: help test test-integration test-unit build install clean docker-build lint fmt vet

# Variables
HOSTNAME=registry.terraform.io
NAMESPACE=EinDev
NAME=snitchdns
BINARY=terraform-provider-${NAME}
VERSION=dev
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

# Default target
help:
	@echo "Available targets:"
	@echo "  test              - Run all tests"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests with testcontainer"
	@echo "  build             - Build the provider"
	@echo "  install           - Install provider locally for development"
	@echo "  clean             - Clean build artifacts"
	@echo "  docker-build      - Build the test container image"
	@echo "  lint              - Run linters"
	@echo "  fmt               - Format code"
	@echo "  vet               - Run go vet"
	@echo "  install-tools     - Install development tools"

# Run all tests
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run unit tests only (fast)
test-unit:
	go test -v -short ./...

# Run integration tests with testcontainer
test-integration:
	go test -v -tags=integration ./internal/testcontainer/...

# Build the provider
build:
	go build -v ./...

# Install the provider locally for development
install: build
	@echo "Installing provider to ${PLUGIN_DIR}..."
	@mkdir -p ${PLUGIN_DIR}
	@go build -o ${PLUGIN_DIR}/${BINARY}_v${VERSION}
	@echo "Provider installed successfully!"
	@echo ""
	@echo "To use the local provider, create a ~/.terraformrc file with:"
	@echo ""
	@echo "provider_installation {"
	@echo "  dev_overrides {"
	@echo "    \"${HOSTNAME}/${NAMESPACE}/${NAME}\" = \"${PLUGIN_DIR}\""
	@echo "  }"
	@echo "  direct {}"
	@echo "}"
	@echo ""
	@echo "Or use the generated .terraformrc file in this directory"

# Clean build artifacts
clean:
	rm -rf ./bin
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}
	rm -f coverage.txt
	go clean -cache -testcache

# Build the test container image manually
docker-build:
	docker build -t snitchdns-test:latest ./testcontainer

# Run linters
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Install development tools
install-tools:
	@echo "Installing golangci-lint v2.7.2..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.7.2
	@echo "Development tools installed successfully"
