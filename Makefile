BINARY_NAME=forgectl
VERSION=0.2.0
BUILD_DIR=dist
LDFLAGS=-ldflags "-s -w -X github.com/forgectl/forgectl/cmd.version=$(VERSION)"

.PHONY: all build clean install uninstall test lint vet dev

all: build

build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

dev:
	@echo "Building for development..."
	go build -o $(BINARY_NAME) .

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/

uninstall:
	@echo "Removing $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

cross:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Built binaries in $(BUILD_DIR)/"

test:
	go test -v ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)