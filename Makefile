BUILD_DIR := build
BINARY := $(BUILD_DIR)/bin
HASURA_GRAPHQL_ADMIN_SECRET := my-admin-secret
HASURA_ENDPOINT_URI := http://localhost:8081
HASURA_METADATA_DIR := hasura

.PHONY: all clean

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build-bin:
	@echo " ## ##  BUILDING   ## ## "
	@mkdir -p $(BUILD_DIR)
	go mod download
	go build -o $(BINARY) main.go

run: build-bin
	@echo " ## ##   RUNNING   ## ## "
	@$(BINARY)

psql:
	docker exec -it go-graphql-trial-postgres-1 psql -U postgres -d userapp

build-dev:
	docker compose -f dev.compose.yml build

up-dev:
	docker compose -f dev.compose.yml up

build-up-dev: build-dev up-dev

down-dev:
	docker compose -f dev.compose.yml down

export-hasura-metadata:
	hasura metadata export --endpoint $(HASURA_ENDPOINT_URI) --admin-secret $(HASURA_GRAPHQL_ADMIN_SECRET) --project $(HASURA_METADATA_DIR)
