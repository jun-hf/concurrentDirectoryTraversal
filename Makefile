build:
	@go build -o bin/td

run: build
	@./bin/td

test:
	@go test ./.. -v
