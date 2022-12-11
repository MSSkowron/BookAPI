BINARY_NAME='gobookapi'

build:
	@go build -o bin/${BINARY_NAME}

docker:
	docker-compose -f docker-compose.yml up -d

run: docker build
	@./bin/${BINARY_NAME}