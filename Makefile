.PHONY: test

all: test

test: ## Test the library
	@mkdir -p artifacts/profiles
	go test ./... -bench=. -race -timeout 10000ms -coverprofile cover.out
	go tool cover -func=cover.out

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
