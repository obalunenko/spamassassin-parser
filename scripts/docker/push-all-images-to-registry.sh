#!/bin/sh

set -eu

SCRIPT_NAME="$(basename "$0")"

echo "${SCRIPT_NAME} is running... "

SPAMASSASSIN_DOCKER_REPO=$1

if ! command -v docker
then
 printf "Cannot check docker, please install docker:
        https://docs.docker.com/get-docker/ \n"
   exit 1
fi

echo "Pushing images to repo ${SPAMASSASSIN_DOCKER_REPO}..."

echo

for i in $(docker images --format "{{.Repository}}:{{.Tag}}" | grep "${SPAMASSASSIN_DOCKER_REPO}" | grep -v "<none>"); do
  echo "[DEBUG]: docker push $i"
  docker push "${i}"
done

echo "${SCRIPT_NAME} done."
