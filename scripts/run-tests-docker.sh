#!/usr/bin/env sh
set -e

export GO111MODULE=on
go test -v -coverpkg=./... -covermode=atomic -coverprofile=coverage.out ./...