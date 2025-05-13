BUILD_DIR := build
BINARY := $(BUILD_DIR)/bin

.PHONY: all clean

fmt:
	@go fmt ./...

vet:
	@go vet ./...

build: fmt vet
	@echo " ## ##  BUILDING   ## ## "
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY) main.go

run: build
	@echo " ## ##   RUNNING   ## ## "
	@$(BINARY)
