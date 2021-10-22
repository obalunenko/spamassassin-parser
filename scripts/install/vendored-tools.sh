#!/bin/sh

set -eu

SCRIPT_NAME="$(basename "$0")"
SCRIPT_DIR="$(dirname "$0")"
REPO_ROOT="$(cd "${SCRIPT_DIR}" && git rev-parse --show-toplevel)"
TOOLS_DIR="${REPO_ROOT}/tools"

echo "${SCRIPT_NAME} is running... "

cd "${TOOLS_DIR}" || exit 1

go generate -x -tags=tools

cd - || exit 1

echo "${SCRIPT_NAME} done."
