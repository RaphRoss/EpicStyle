.PHONY: build test clean install

BINARY_NAME=epicstyle
BUILD_DIR=build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

lint:
	golangci-lint run

format:
	go fmt ./...