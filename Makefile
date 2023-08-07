BINARY_NAME=bookrestapi

build:
	@go build -o bin/$(BINARY_NAME) ./cmd/bookrestapi/main.go

run: build
	@./bin/$(BINARY_NAME)

test: 
	@go test -race -vet=off ./...

lint:
	staticcheck ./...
	golint ./...

clean:
	@rm -rf bin/$(BINARY_NAME)

.PHONY: build run test lint clean