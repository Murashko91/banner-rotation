BIN := "./bin/banners_rotation"
DB_NAME := "banners_rotation"
DB_CONTAINER_NAME := "banners_rotation_psql"

DOCKER_IMG="banner-app:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/.
	
run: swag-init build 
	$(BIN) -conf ./configs/config.yaml

db-up:
	docker run --name $(DB_CONTAINER_NAME) -p 5454:5432 -e POSTGRES_USER=otus -e POSTGRES_PASSWORD=otus -e POSTGRES_DB=$(DB_NAME) -d -v "./migrations":/docker-entrypoint-initdb.d postgres:13.3
db-down:
	docker stop $(DB_CONTAINER_NAME)
	docker remove $(DB_CONTAINER_NAME)

swag-init:
	swag init -d ./cmd/,./internal/server/http/,./internal/storage -g main.go --parseDependency

build-img: swag-init
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run  -p 8888:8888 $(DOCKER_IMG) 

version: build
	$(BIN) version

test:
	go test -race -count 100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.3

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker compose -f docker-compose.yaml up -d 

up-rebuild: build-img up

down:
	docker compose -f docker-compose.yaml down 

push: build-img
	docker tag banner-app:develop murashkosv91/banner-app:develop
	docker push murashkosv91/banner-app:develop

integration-tests:
	go test -v ./test/integration/integration_test.go

.PHONY: build run build-img run-img version test lint swag-init integration-tests push
