.PHONY: build test clean install lint format fmtcheck check install-tools

BINARY_NAME=bm
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

# Check if GOPATH is set, if not use ~/go
ifeq ($(GOPATH),)
GOPATH=$(HOME)/go
$(info GOPATH not set, using default: $(GOPATH))
endif

install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	@go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/bm

test:
	@go test -v ./...

lint:
	@echo "Running linter..."
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run

formatcheck:
	@echo "Checking formatting..."
	@if [ -n "$(shell gofmt -l .)" ]; then \
		echo "Wrong formatting for files:"; \
		gofmt -l .; \
		echo "Run 'make format' to fix formatting issues."; \
		exit 1; \
	fi
	@go mod tidy


format:
	@echo "Formatting code..."
	@gofmt -s -w .
	@go mod tidy

# Combined check target that runs format, lint, and test
check: format lint test
	@echo "All checks passed!"

clean:
	@rm -rf ${BUILD_DIR}

install: build
	@echo "Installing ${BINARY_NAME}..."
	@cp ${BUILD_DIR}/${BINARY_NAME} ${GOPATH}/bin/
	@echo "Installation complete"
