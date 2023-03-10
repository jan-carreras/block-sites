#!/usr/bin/env make

.PHONY: help
help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Pass all the tests
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

build: test ## Build the binary
	go build -o bin/hosts

.PHONY: help
install: build ## Install the binary in the system
	cp bin/hosts ~/.bin/hosts
