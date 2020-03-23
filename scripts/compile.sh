#!/usr/bin/env sh
set -e
echo "Building..."

BIN_OUT=./bin/spamassassin-parser-cli

go build -o ${BIN_OUT} ./cmd/spamassassin-parser-cli

echo "Binary compiled at ${BIN_OUT}"