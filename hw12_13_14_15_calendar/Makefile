CALENDAR_BIN := "./bin/calendar"
CALENDAR_DOCKER_IMG="calendar:develop"
DOCKER_COMPOSE_FOLDER="./deployments"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
.PHONY: help

up:
	cd $(DOCKER_COMPOSE_FOLDER) && docker-compose up

rebuild:
	cd $(DOCKER_COMPOSE_FOLDER) && docker-compose up --build

down:
	cd $(DOCKER_COMPOSE_FOLDER) && docker-compose down --remove-orphans

build-calendar:
	go build -v -o $(CALENDAR_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: run-calendar

run-calendar: build-calendar
	$(CALENDAR_BIN) -config ./configs/config.toml

build: build-calendar

build-img: build-calendar-img

build-calendar-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(CALENDAR_DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-calendar-img
	docker run $(CALENDAR_DOCKER_IMG)

version: build-calendar
	$(CALENDAR_BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

pre-commit: ## Pre commit handles
	go mod tidy
	go mod verify
	gofmt -w -s -d .
	golangci-lint run ./...
.PHONY: pre-commit

generate: ## Run go generate
	go generate ./...

test-packages: ## Test all packages
	go test all

wire: ## Run DI generating
	wire cmd/calendar/wire.go

.PHONY: build run build-img run-img version test lint
