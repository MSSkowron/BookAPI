BINARY_NAME='gobankapi'

build:
	@go build -o bin/${BINARY_NAME}

run: build 
	@./bin/${BINARY_NAME}