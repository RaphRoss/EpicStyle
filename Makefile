# Makefile pour compiler main.go

BINARY_NAME=epicstyle

all: build

build:
	go build -o $(BINARY_NAME) main.go

clean:
	rm -f $(BINARY_NAME)

.PHONY: all build clean