.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --no-builtin-rules

.PHONY: test vet install lint build

.DEFAULT_GOAL := build

vet: ## run go vet
	go vet ./...

build: # build opms binary
	@echo " > building opms binary"
	@go build -o opms .
	@echo " - build complete"

test:
	go test -timeout 1m ./...

lint:
	golangci-lint run --fix

install: ## install required dependencies
	@echo "> installing dependencies"
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
