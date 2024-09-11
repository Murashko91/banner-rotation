BIN := "./bin/banners_rotation"
DB_NAME := "banners_rotation"
DB_CONTAINER_NAME := "banners_rotation_psql"

DOCKER_IMG="banners_rotation:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/.
	
run: build
	$(BIN) -conf ./configs/config.yaml

db-up:
	docker run --name $(DB_CONTAINER_NAME) -p 5454:5432 -e POSTGRES_USER=otus -e POSTGRES_PASSWORD=otus -e POSTGRES_DB=$(DB_NAME) -d -v "./migrations":/docker-entrypoint-initdb.d postgres:13.3
db-down:
	docker stop $(DB_CONTAINER_NAME)
	docker remove $(DB_CONTAINER_NAME)
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
	go test -race -count 100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.3

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
