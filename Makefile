.PHONY: test

all: test

clean:
	@rm -rf artifacts

test: clean ## Test the library
	@mkdir -p artifacts
	@go list ./... | grep -v vendor | xargs go vet
	@go list ./... | grep -v vendor | xargs golint
	@go test ./... -bench=. -race -timeout 10000ms -coverprofile cover.out
	@go tool cover -html=cover.out -o artifacts/cover.html
	@go tool cover -func=cover.out

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
