help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
.PHONY: help

pre-commit: ## Pre commit handles
	go mod tidy
	gofumpt -l -w .
	golangci-lint run ./...
.PHONY: pre-commit

test: ## Run tests with cover
	go test ./... -cover -v
.PHONY: test

race-test: ## Run tests with cover and race
	go test ./... -race -v
.PHONY: race-test

bench-old: ## Run old bench
	go test -bench=. -benchmem -count 5 | tee old

bench-new: ## Run new bench
	go test -bench=. -benchmem -count 5 | tee new

profile-export:
	go test -bench=Benchmar -cpuprofile=cpu.out -memprofile=mem.out

pprof-mem: profile-export ## Run pprof for mem
	go tool pprof -http=":8090" mem.out

pprof-cpu: profile-export ## Run pprof for cpu
	go tool pprof -http=":8090" mem.out

bench-stat: ## Run bench stat
	benchstat -delta-test none old new

task-performance-test: ## test-performance
	go test -v -count=1 -timeout=30s -tags bench .

pipe: pre-commit race-test task-performance-test ## Run pre commit tasks

performance: bench-new bench-stat ## Run a compare

.DEFAULT_GOAL := help
