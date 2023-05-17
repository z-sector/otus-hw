BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DCF := -f deployments/docker-compose.yaml --env-file deployments/.env.example

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	APP_STORAGE_DB_USERNAME=user APP_STORAGE_DB_PASSWORD=1234 $(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	#go test -race ./internal/... ./pkg/...
	go test -v -race -count=1 ./...

install-lint-deps:
	@(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1
#	@(which golangci-lint > /dev/null) || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

dc-down:
	docker compose $(DCF) down

dc-up: dc-down
	docker compose $(DCF) up --build --detach

dc-down-prune:
	docker compose $(DCF) down -v --rmi all

makemigrations: ## e.g `make makemigrations name=init`
	docker compose $(DCF) run --rm --no-deps migrator migrate create -ext sql -dir /migrations $(name)

.PHONY: build run build-img run-img version test lint dc-up dc-down dc-down-prune makemigrations
