NAME=spamassassin-parser-cli
BIN_DIR=./bin

SHELL := env DOCKER_REPO=$(DOCKER_REPO) $(SHELL)
DOCKER_REPO?=olegbalunenko

SHELL := env VERSION=$(VERSION) $(SHELL)
VERSION ?= $(shell git describe --tags $(git rev-list --tags --max-count=1))

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



build: compile-spamassassin-parser-be
.PHONY: build

compile-spamassassin-parser-be:
	./scripts/build/spamassassin-parser-be.sh
.PHONY: compile-spamassassin-parser-be

## Test coverage report.
test-cover:
	${call colored, test-cover is running...}
	./scripts/tests/coverage.sh
.PHONY: test-cover

## Open coverage report.
open-cover-report: test-cover
	./scripts/open-coverage-report.sh
.PHONY: open-cover-report

update-readme-cover: compile test-cover
	./scripts/update-readme-coverage.sh
.PHONY: update-readme-cover

test:
	./scripts/tests/run.sh
.PHONY: test

coverage:
	make cover

configure: sync-vendor

sync-vendor:
	./scripts/sync-vendor.sh
.PHONY: sync-vendor

## Fix imports sorting.
imports:
	${call colored, fix-imports is running...}
	./scripts/style/fix-imports.sh
.PHONY: imports

## Format code with go fmt.
fmt:
	${call colored, fmt is running...}
	./scripts/style/fmt.sh
.PHONY: fmt

## Format code and sort imports.
format-project: fmt imports
.PHONY: format-project

install-tools:
	./scripts/install/vendored-tools.sh
.PHONY: install-tools

## vet project
vet:
	${call colored, vet is running...}
	./scripts/linting/run-vet.sh
.PHONY: vet

## Run full linting
lint-full:
	./scripts/linting/run-linters.sh
.PHONY: lint-full

## Run linting for build pipeline
lint-pipeline:
	./scripts/linting/run-linters-pipeline.sh
.PHONY: lint-pipeline

## recreate all generated code and swagger documentation.
codegen:
	${call colored, codegen is running...}
	./scripts/codegen/go-generate.sh
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
	${call colored, release is running...}
	./scripts/release/local-snapshot-release.sh
.PHONY: release-local-snapshot

## Issue new release.
new-version: vet test compile
	./scripts/release/new-version.sh
.PHONY: new-release



######################################
############### DOCKER ###############
######################################

################ PROD #################

## Push all prod images to registry.
docker-push-prod-images:
	${call colored, docker-push-prod-images is running...}
	./scripts/docker/push-all-images-to-registry.sh ${DOCKER_REPO}
.PHONY: docker-push-prod-images

## Build docker base images.
docker-build-base-prod: docker-build-base-go-prod
.PHONY: docker-build-base-prod

## Build docker base image for GO
docker-build-base-go-prod:
	${call colored, docker-build-base-go-prod is running...}
	./scripts/docker/build/prod/go-base.sh
.PHONY: docker-build-base-go-prod

## Build all services docker prod images for deploying to gcloud.
docker-build-prod: docker-build-backend-prod
.PHONY: docker-build-prod

## Build all backend services docker prod images for deploying to gcloud.
docker-build-backend-prod: docker-build-spamassassin-parser-prod
.PHONY: docker-build-backend-prod

## Build admin service prod docker image.
docker-build-spamassassin-parser-prod:
	${call colored, docker-build-prod backend-admin is running...}
	./scripts/docker/build/prod/spamassassin-parser.sh
.PHONY: docker-build-spamassassin-parser-prod

## Docker compose up - deploys prod containers on docker locally.
docker-compose-up:
	${call colored, docker-up is running...}
	./scripts/docker/compose/prod/up.sh
.PHONY: docker-compose-up

## Docker compose down - remove all prod containers in docker locally.
docker-compose-down:
	${call colored, docker-down is running...}
	./scripts/docker/compose/prod/down.sh
.PHONY: docker-compose-down

## Docker compose stop - stops all prod containers in docker locally.
docker-compose-stop:
	${call colored, docker-down is running...}
	./scripts/docker/compose/prod/stop.sh
.PHONY: docker-compose-stop

## Build all prod images: base and services.
docker-prepare-images-prod: docker-build-base-prod docker-build-prod
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

## Open containers logs service url.
open-container-logs:
	./scripts/browser-opener.sh -u 'http://localhost:9999/'
.PHONY: open-container-logs


################## DEV ###################

## Build docker base images.
docker-build-base-dev: docker-build-base-go-dev
.PHONY: docker-build-base-dev

## Build docker base image for GO
docker-build-base-go-dev:
	${call colored, docker-build-base-go-dev is running...}
	./scripts/docker/build/dev/go-base.sh
.PHONY: docker-build-base-go-dev

## Build docker dev image for running locally.
docker-build-dev: docker-build-spamassassin-parser-dev
.PHONY: docker-build-dev

## Build admin service dev docker image.
docker-build-spamassassin-parser-dev:
	${call colored, docker-build-admin-dev is running...}
	./scripts/docker/build/dev/spamassassin-parser.sh
.PHONY: docker-build-spamassassin-parser-dev

## Dev Docker-compose up with stubbed 3rd party dependencies.
dev-docker-compose-up:
	${call colored, dev-docker-up is running...}
	./scripts/docker/compose/dev/up.sh
.PHONY: dev-docker-compose-up

## Docker compose down.
dev-docker-compose-down:
	${call colored, dev-docker-down is running...}
	./scripts/docker/compose/dev/down.sh
.PHONY: dev-docker-compose-down

## Docker compose stop - stops all dev containers in docker locally.
dev-docker-compose-stop:
	${call colored, docker-down is running...}
	./scripts/docker/compose/dev/stop.sh
.PHONY: dev-docker-compose-stop

## Dev local full deploy: build base images, build services images, deploy to docker compose
deploy-local-dev: docker-build-base-dev docker-build-dev run-local-dev
.PHONY: deploy-local-dev

## Run locally dev: deploy to docker compose and expose tunnels.
run-local-dev: dev-docker-compose-up
.PHONY: run-local-dev

## Stop the world and close tunnels.
stop-local-dev: dev-docker-compose-stop
.PHONY: stop-local-prod

.DEFAULT_GOAL := help

