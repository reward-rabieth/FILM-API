
build:
	@go build -o bin/gofilm

run: build
	@./bin/gofilm

test:
	@go test -v ./...