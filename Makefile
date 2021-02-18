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



## Compile executable
compile:
	./scripts/compile.sh
.PHONY: compile

## lint project
lint:
	./scripts/run-linters.sh
.PHONY: lint

lint-ci:
	./scripts/run-linters-ci.sh
.PHONY: lint-ci


## format markdown files in project
pretty-markdown:
	find . -name '*.md' -not -wholename './vendor/*' | xargs prettier --write
.PHONY: pretty-markdown

## Test all packages
test:
	./scripts/run-tests.sh
.PHONY: test

## Test coverage
test-cover:
	./scripts/coverage.sh
.PHONY: test-cover

new-version: lint test compile
	./scripts/version.sh
.PHONY: new-version

## Release
release:
	./scripts/release.sh
.PHONY: release

## Release local snapshot
release-local-snapshot:
	${call colored, release is running...}
	./scripts/local-snapshot-release.sh
.PHONY: release-local-snapshot

## Fix imports sorting.
imports:
	${call colored, fix-imports is running...}
	./scripts/fix-imports.sh
.PHONY: imports

## Format code.
fmt:
	${call colored, fmt is running...}
	./scripts/fmt.sh
.PHONY: fmt

## Format code and sort imports.
format-project: fmt imports
.PHONY: format-project

## fetch all dependencies for scripts
install-tools:
	./scripts/get-dependencies.sh
.PHONY: install-tools

## Sync vendor
sync-vendor:
	${call colored, gomod is running...}
	./scripts/sync-vendor.sh
.PHONY: sync-vendor

## Update dependencies
gomod-update:
	${call colored, gomod is running...}
	go get -u -v ./...
	make sync-vendor
.PHONY: gomod-update

vet:
	./scripts/vet.sh
.PHONY: vet

## Docker compose up
docker-compose-up:
	docker-compose -f ./docker-compose.yml up --build -d

.PHONY: docker-compose-up

## Docker compose down
docker-compose-down:
	docker-compose -f ./docker-compose.yml down --volumes

.PHONY: docker-compose-down

## Docker compose up
docker-compose-up-dev:
	docker-compose -f ./docker-compose.dev.yml up --build -d

.PHONY: docker-compose-up-dev

## Docker compose down
docker-compose-down-dev:
	docker-compose -f ./dev.docker-compose.dev.yml down --volumes

.PHONY: docker-compose-down-dev

## Build docker base image for GO
docker-build-base-go-prod:
	${call colored, docker-build-base-go-prod is running...}
	docker build --rm --no-cache -t ${DOCKER_REPO}/spamassassin-go-base:${VERSION} -t ${DOCKER_REPO}/spamassassin-go-base:latest -f ./build/docker/base-docker/go.Dockerfile .
.PHONY: docker-build-base-go-prod

## Build admin service prod docker image.
docker-build-spamassassin-prod:
	${call colored, docker-build-spamassassin-prod is running...}
	docker build --rm --no-cache -t ${DOCKER_REPO}/spamassassin-parser:${VERSION} -t ${DOCKER_REPO}/spamassassin-parser:latest -f ./build/docker/spamassassin-parser/Dockerfile .
.PHONY: docker-build-admin-prod

docker-build-prod: docker-build-base-go-prod docker-build-spamassassin-prod
.PHONY: docker-build-prod

.DEFAULT_GOAL := test

