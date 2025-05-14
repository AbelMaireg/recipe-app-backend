BUILD_DIR := build
BINARY := $(BUILD_DIR)/bin

.PHONY: all clean

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build:
	@echo " ## ##  BUILDING   ## ## "
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) main.go

run: build
	@echo " ## ##   RUNNING   ## ## "
	@$(BINARY)

psql:
	docker exec -it go-graphql-trial-postgres-1 psql -U postgres -d userapp

gqlgen-init:
	go run github.com/99designs/gqlgen init

gqlgen-generate:
	go run github.com/99designs/gqlgen generate

build-dev:
	docker compose -f dev.compose.yml build

up-dev:
	docker compose -f dev.compose.yml up

build-up-dev: build-dev up-dev

down-dev:
	docker compose -f dev.compose.yml down
