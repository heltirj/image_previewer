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

run: build
	./$(BINARY_NAME)

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)

deps:
	go mod tidy

.PHONY: all build run test clean deps