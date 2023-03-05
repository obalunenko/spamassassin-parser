#!/bin/bash

set -Eeuo pipefail

SCRIPT_NAME="$(basename "$0")"
SCRIPT_DIR="$(dirname "$0")"
REPO_ROOT="$(cd "${SCRIPT_DIR}" && git rev-parse --show-toplevel)"
SCRIPTS_DIR="${REPO_ROOT}/scripts"

source "${SCRIPTS_DIR}/helpers-source.sh"

echo "${SCRIPT_NAME} is running... "

URL=""

while getopts u: flag; do
  case "${flag}" in
  u)
    URL=${OPTARG}
    ;;
  *)
    echo "Unknown flag passed"
    exit 1
    ;;
  esac
done

function openurl() {
  echo "${URL}"
  openSource "${URL}"
}

openurl

echo "${SCRIPT_NAME} done."
