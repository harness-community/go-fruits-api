current_dir = $(shell pwd)
compose_file = "$(current_dir)/docker-compose.yaml"
FRUITS_DB_SERVICE ?= "postgresql"
TEST_LOG_LEVEL ?= "info"

swaggo:	## Generate Swagger OpenAPI docs
	@swag  init --parseDependency --parseInternal -g server.go

test:	## Runs test
	@drone exec --trusted --env-file=.env --include=test --include=$(FRUITS_DB_SERVICE)

start-db:	## Starts the docker containers using docker-compose
	@docker-compose -f $(compose_file) up -d

clean:	## Cleans output
	go clean
	rm -rf dist

vendor:	## Vendoring
	go mod vendor

lint:	## Run lint on the project
	@drone exec --trusted --include=lint --env-file=.env

ko: ## Dev deployment using ko
	kustomize build config/app | ko resolve --platform=linux/arm64 -f - | kubectl apply -f -

manifests:	## Generates application deployment manifests

build-push-image:	## Builds Container Image
	@drone exec --trusted --env-file=.env

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'

.PHONY: swaggo	test	start-db	clean	lint	vendor	help swaggo build-push-image