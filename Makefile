build:
	@go build -o bin/githubAssistant

run: build
	@./bin/githubAssistant

test:
	@go test -v ./...
