build:
	@go build -o bin/go-matcher cmd/main.go

run: build
	@./bin/go-matcher

test-betting:
	@go test -v ./services/betting