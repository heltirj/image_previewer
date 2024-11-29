# Укажите имя вашего исполняемого файла
BINARY_NAME=./bin/image_previewer

# Укажите вашу основную директорию
SRC_DIR=./cmd/api

# Укажите возможности для компиляции (различные архитектуры)
GOOS=linux
GOARCH=amd64

all: build

build:
	go build -o $(BINARY_NAME) $(SRC_DIR)/*.go

run:
	@docker-compose up --build

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)

deps:
	go mod tidy

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: all build run test clean deps install-lint-deps lint