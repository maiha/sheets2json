# Get version info from git
GIT_DESC := $(shell git describe --tags --long --always --dirty 2>/dev/null || echo "unknown")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -I)

# Process version
ifeq ($(findstring -0-g,$(GIT_DESC)),-0-g)
	VERSION_NUM := $(shell echo $(GIT_DESC) | sed 's/-0-g.*//' | sed 's/^v//')
else ifneq ($(findstring -g,$(GIT_DESC)),)
	VERSION_NUM := $(shell echo $(GIT_DESC) | sed -E 's/^v?([^-]*)-([0-9]+)-g.*/\1+\2/')
else
	VERSION_NUM := $(GIT_DESC)
endif

# Complete version string
VERSION_STRING := sheets2json $(VERSION_NUM) [$(GIT_COMMIT)] ($(BUILD_DATE))

.PHONY: build clean version help fmt lint test

# Default target
build: version
	@echo "Building sheets2json..."
	@echo "$(VERSION_STRING)" > version.txt
	@docker run --rm \
		-v "$${PWD}":/src \
		-v go-mod-cache:/go/pkg/mod \
		-w /src \
		golang:1.21-alpine sh -c "\
		go mod tidy && \
		go build -ldflags=\"-s -w \
			-X 'main.versionString=\$$(cat version.txt)'\" \
		-o sheets2json main.go && \
		chown $$(id -u):$$(id -g) sheets2json"
	@echo "Binary created: ./sheets2json"

# Show version info
version:
	@echo "$(VERSION_STRING)" > version.txt
	@echo "$(VERSION_STRING)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f sheets2json version.txt

# Format Go code
fmt:
	@docker run --rm \
		-u "$$(id -u):$$(id -g)" \
		-e HOME=/tmp \
		-v "$${PWD}":/src \
		-v go-mod-cache:/go/pkg/mod \
		-w /src \
		golang:1.21-alpine go fmt ./...

# Lint Go code
lint:
	@docker run --rm \
		-u "$$(id -u):$$(id -g)" \
		-e HOME=/tmp \
		-v "$${PWD}":/src \
		-v go-mod-cache:/go/pkg/mod \
		-w /src \
		golangci/golangci-lint:v1.54-alpine golangci-lint run

# Run tests
test:
	@docker run --rm \
		-v "$${PWD}":/src \
		-v go-mod-cache:/go/pkg/mod \
		-w /src \
		golang:1.21-alpine go test -v ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  make build     - Build binary (default)"
	@echo "  make fmt       - Format Go code"
	@echo "  make lint      - Lint Go code"
	@echo "  make test      - Run tests"
	@echo "  make version   - Show version information"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make help      - Show this help message"
