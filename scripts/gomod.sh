#!/usr/bin/env sh

go mod tidy -v
go mod vendor
go mod verify