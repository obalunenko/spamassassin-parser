#!/usr/bin/env sh
set -e
echo "Building..."

BIN_OUT=./bin/spamassassin-parser

go build -o ${BIN_OUT} ./cmd/spamassassin-parser

echo "Binary compiled at ${BIN_OUT}"