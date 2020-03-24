#!/usr/bin/env bash

# Get new tags from the remote
git fetch --tags

# Get the latest tag name
latestTag=$(git describe --tags $(git rev-list --tags --max-count=1))
echo "${latestTag}"

echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

curl -sL https://git.io/goreleaser | bash