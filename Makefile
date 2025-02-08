.PHONY: build test clean install

BINARY_NAME=bookmarker
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

# Check if GOPATH is set, if not use ~/go
ifeq ($(GOPATH),)
GOPATH=$(HOME)/go
$(info GOPATH not set, using default: $(GOPATH))
endif

build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	@go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/bookmarker

test:
	@go test -v ./...

clean:
	@rm -rf ${BUILD_DIR}

install: build
	@echo "Installing ${BINARY_NAME}..."
	@cp ${BUILD_DIR}/${BINARY_NAME} ${GOPATH}/bin/
	# @mkdir -p ~/.config/fish/completions
	# @cp pkg/shell/completions/bm.fish ~/.config/fish/completions/
	# @cp pkg/shell/wrapper.fish ~/.config/fish/functions/bm.fish
	@echo "Installation complete"
