BIN_DIR=./bin

DOCKER_REPO ?= ghcr.io/obalunenko/
export DOCKER_REPO

SHELL := env VERSION=$(VERSION) $(SHELL)
VERSION ?= $(shell git describe --tags $(git rev-list --tags --max-count=1))

APP_NAME?=spamassassin-parser
SHELL := env APP_NAME=$(APP_NAME) $(SHELL)

GOTOOLS_IMAGE_TAG?=v0.6.1
SHELL := env GOTOOLS_IMAGE_TAG=$(GOTOOLS_IMAGE_TAG) $(SHELL)

COMPOSE_TOOLS_FILE=deployments/docker-compose/go-tools-docker-compose.yml
COMPOSE_TOOLS_CMD_BASE=docker compose -f $(COMPOSE_TOOLS_FILE)
COMPOSE_TOOLS_CMD_UP=$(COMPOSE_TOOLS_CMD_BASE) up --exit-code-from
COMPOSE_TOOLS_CMD_PULL=$(COMPOSE_TOOLS_CMD_BASE) pull

TARGET_MAX_CHAR_NUM=20

## Show help
help:
	${call colored, help is running...}
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-$(TARGET_MAX_CHAR_NUM)s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)



## Build project.
build: compile-app
.PHONY: build

## Compile app.
compile-app:
	./scripts/build/app.sh
.PHONY: compile-app

## Test coverage report.
test-cover:
	./scripts/tests/coverage.sh
.PHONY: test-cover

prepare-cover-report: test-cover
	$(COMPOSE_TOOLS_CMD_UP) prepare-cover-report prepare-cover-report
.PHONY: prepare-cover-report

## Open coverage report.
open-cover-report: prepare-cover-report
	./scripts/open-coverage-report.sh
.PHONY: open-cover-report

## Update readme coverage.
update-readme-cover: build prepare-cover-report
	$(COMPOSE_TOOLS_CMD_UP) update-readme-coverage update-readme-coverage
.PHONY: update-readme-cover

## Run tests.
test:
	$(COMPOSE_TOOLS_CMD_UP) run-tests run-tests
.PHONY: test

## Run regression tests.
test-regression: test
.PHONY: test-regression

## Sync vendor and install needed tools.
configure: sync-vendor install-tools

## Sync vendor with go.mod.
sync-vendor:
	./scripts/sync-vendor.sh
.PHONY: sync-vendor

## Fix imports sorting.
imports:
	$(COMPOSE_TOOLS_CMD_UP) fix-imports fix-imports
.PHONY: imports

## Format code with go fmt.
fmt:
	$(COMPOSE_TOOLS_CMD_UP) fix-fmt fix-fmt
.PHONY: fmt

## Format code and sort imports.
format-project: fmt imports
.PHONY: format-project

## Installs vendored tools.
install-tools:
	echo "Installing ${GOTOOLS_IMAGE_TAG}"
	$(COMPOSE_TOOLS_CMD_PULL)
.PHONY: install-tools

## vet project
vet:
	./scripts/linting/run-vet.sh
.PHONY: vet

## Run full linting
lint-full:
	$(COMPOSE_TOOLS_CMD_UP) lint-full lint-full
.PHONY: lint-full

## Run linting for build pipeline
lint-pipeline:
	$(COMPOSE_TOOLS_CMD_UP) lint-pipeline lint-pipeline
.PHONY: lint-pipeline

## Run linting for sonar report
lint-sonar:
	$(COMPOSE_TOOLS_CMD_UP) lint-sonar lint-sonar
.PHONY: lint-sonar

## recreate all generated code and documentation.
codegen:
	$(COMPOSE_TOOLS_CMD_UP) go-generate go-generate
.PHONY: codegen

## recreate all generated code and swagger documentation and format code.
generate: codegen format-project vet
.PHONY: generate

## Release
release:
	./scripts/release/release.sh
.PHONY: release

## Release local snapshot
release-local-snapshot:
	./scripts/release/local-snapshot-release.sh
.PHONY: release-local-snapshot

## Check goreleaser config.
check-releaser:
	./scripts/release/check.sh
.PHONY: check-releaser

## Issue new release.
new-version: vet test-regression build
	./scripts/release/new-version.sh
.PHONY: new-release



######################################
############### DOCKER ###############
######################################


################ BASE #################
## Build docker base images.
docker-build-base: docker-build-base-alpine docker-build-base-go
.PHONY: docker-build-base

## Build docker base image for GO
docker-build-base-go:
	./scripts/docker/build/base/go.sh
.PHONY: docker-build-base-go

## Build docker base image for GO
docker-build-base-alpine:
	./scripts/docker/build/base/alpine.sh
.PHONY: docker-build-base-alpine


################ PROD #################

## Push all prod images to registry.
docker-push-prod-images:
	./scripts/docker/push-all-images-to-registry.sh ${DOCKER_REPO}
.PHONY: docker-push-prod-images

## Build all services docker prod images for deploying to gcloud.
docker-build-prod: docker-build-backend-prod
.PHONY: docker-build-prod

## Build all backend services docker prod images for deploying to gcloud.
docker-build-backend-prod: docker-build-spamassassin-prod
.PHONY: docker-build-backend-prod

## Build admin service prod docker image.
docker-build-spamassassin-prod:
	./scripts/docker/build/prod/spamassassin-parser.sh
.PHONY: docker-build-spamassassin-prod

## Docker compose up - deploys prod containers on docker locally.
docker-compose-up:
	./scripts/docker/compose/prod/up.sh
.PHONY: docker-compose-up

## Docker compose down - remove all prod containers in docker locally.
docker-compose-down:
	./scripts/docker/compose/prod/down.sh
.PHONY: docker-compose-down

## Docker compose stop - stops all prod containers in docker locally.
docker-compose-stop:
	./scripts/docker/compose/prod/stop.sh
.PHONY: docker-compose-stop

## Build all prod images: base and services.
docker-prepare-images-prod: docker-build-base docker-build-prod
.PHONY: docker-prepare-images-prod

## Prod local full deploy: build base images, build services images, deploy to docker compose
deploy-local-prod: docker-prepare-images-prod run-local-prod
.PHONY: deploy-local-prod

## Run locally: deploy to docker compose and expose tunnels.
run-local-prod: docker-compose-up
.PHONY: run-local-prod

## Stop the world and close tunnels.
stop-local-prod: docker-compose-stop
.PHONY: stop-local-prod

################## DEV ###################

## Build docker dev image for running locally.
docker-build-dev: docker-build-spamassassin-dev
.PHONY: docker-build-dev

## Build admin service dev docker image.
docker-build-spamassassin-dev:
	./scripts/docker/build/dev/spamassassin-parser.sh
.PHONY: docker-build-spamassassin-dev

## Dev Docker-compose up with stubbed 3rd party dependencies.
dev-docker-compose-up:
	./scripts/docker/compose/dev/up.sh
.PHONY: dev-docker-compose-up

## Docker compose down.
dev-docker-compose-down:
	./scripts/docker/compose/dev/down.sh
.PHONY: dev-docker-compose-down

## Docker compose stop - stops all dev containers in docker locally.
dev-docker-compose-stop:
	./scripts/docker/compose/dev/stop.sh
.PHONY: dev-docker-compose-stop

## Build all dev images: base and services.
docker-prepare-images-dev: docker-build-base docker-build-dev
.PHONY: docker-prepare-images-dev

## Dev local full deploy: build base images, build services images, deploy to docker compose
deploy-local-dev: docker-prepare-images-dev run-local-dev
.PHONY: deploy-local-dev

## Run locally dev: deploy to docker compose and expose tunnels.
run-local-dev: dev-docker-compose-up
.PHONY: run-local-dev

## Stop the world and close tunnels.
stop-local-dev: dev-docker-compose-stop
.PHONY: stop-local-prod

## Open containers logs service url.
open-container-logs:
	./scripts/browser-opener.sh -u 'http://localhost:9999/'
.PHONY: open-container-logs

.DEFAULT_GOAL := help

