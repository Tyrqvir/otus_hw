CALENDAR_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/calendar_scheduler"
SENDER_BIN := "./bin/calendar_sender"
CALENDAR_DOCKER_IMG :="calendar:develop"
SCHEDULER_DOCKER_IMG :="scheduler:develop"
SENDER_DOCKER_IMG :="sender:develop"
INTERGRATION_TEST_DOCKER_IMG :="integration-test:develop"
DOCKER_COMPOSE_FOLDER := ./deployments
DOCKER_FILE := $(DOCKER_COMPOSE_FOLDER)/docker-compose.yaml
DOCKER_TEST_FILE := $(DOCKER_COMPOSE_FOLDER)/docker-compose.test.yaml
DOCKER_TEST_ENV := $(DOCKER_COMPOSE_FOLDER)/.env.testing

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'
.PHONY: help

up-watch: ## docker up with watcher
	cd $(DOCKER_COMPOSE_FOLDER) && docker-compose up

up: ## docker up
	cd $(DOCKER_COMPOSE_FOLDER) && docker-compose up -d

rebuild: ## docker rebuild
	cd $(DOCKER_COMPOSE_FOLDER) && docker-compose up -d --build

down: ## docker down
	cd $(DOCKER_COMPOSE_FOLDER) && docker-compose down --remove-orphans

build-calendar: ## build calendar
	go build -v -o $(CALENDAR_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: run-calendar## run all binaries

run-calendar: build-calendar ## run calendar
	$(CALENDAR_BIN) -config ./configs/calendar_config.toml

build: build-calendar ## build all app

build-img: build-calendar-img ## build all images

build-calendar-img: ## build calendar image
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--build-arg=BIN_FILE="/opt/calendar/calendar-app" \
		--build-arg=APP_FOLDER="calendar" \
		--build-arg=CONFIG_FILE_EXTERNAL="configs/calendar_config.toml" \
		-t $(CALENDAR_DOCKER_IMG) \
		-f build/Dockerfile .

run-calendar-img: build-calendar-img ## run calendar image
	docker run $(CALENDAR_DOCKER_IMG)

run-img: run-calendar-img ## run all images

version: build-calendar ## show version of calendar
	$(CALENDAR_BIN) version

test: ## run test
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps ## run linter
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

migrate:
	goose -dir migrations/ postgres "host=localhost port=5432 user=postgres password=megapass dbname=calendar sslmode=disable" up

full: lint tests build run build-img run-img version

build-scheduler: ## build scheduler
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

build-sender: ## build sender
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender


run: run-calendar run-scheduler run-sender ## run all binaries
