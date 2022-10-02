current_dir = $(shell pwd)
compose_file = "$(current_dir)/docker-compose.yaml"
FRUITS_DB_TYPE ?= "pg"
TEST_LOG_LEVEL ?= "info"

swaggo:	## Generate Swagger OpenAPI docs
	@swag  init --parseDependency --parseInternal -g server.go

test:	start-db	## Runs test
	go clean -testcache
	go test ./... -v
	@docker-compose -f $(compose_file) down

start-db:	## Starts the docker containers usig docker-compose
	@docker-compose -f $(compose_file) up -d

clean:	## Cleans output
	go clean
	rm -rf dist

vendor:	## Vendoring
	go mod vendor

lint:	## Run lint on the project
	@golangci-lint run

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'

.PHONY: swaggo	test	start-db	clean	lint	vendor	help
