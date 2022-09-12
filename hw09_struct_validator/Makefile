help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
.PHONY: help

pre-commit: ## Pre commit handles
	go mod tidy
	golangci-lint run ./...
	gofumpt -l -w .
.PHONY: pre-commit

test: ## Run tests with cover
	go test ./... -cover -v
.PHONY: test

race-test: ## Run tests with cover and race
	go test ./... -race -v
.PHONY: ravr-test

pipe: pre-commit race-test ## Run all tasks

.DEFAULT_GOAL := help
