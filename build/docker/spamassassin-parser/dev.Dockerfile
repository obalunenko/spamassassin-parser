ARG DOCKER_REPO_BASE
ARG DOCKER_GO_BASE_DEV_TAG=latest
# hadolint ignore=DL3007
FROM ${DOCKER_REPO_BASE}spamassassin-go-base-dev:${DOCKER_GO_BASE_DEV_TAG} as build-container
LABEL maintainer="oleg.balunenko@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/obalunenko/spamassassin-parser"
LABEL stage="dev"

ENV PROJECT_DIR="${GOPATH}/src/github.com/obalunenko/spamassassin-parser"

RUN mkdir -p "${PROJECT_DIR}"

WORKDIR "${PROJECT_DIR}"

COPY .git .git
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY scripts scripts
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile

# compile executable
RUN make compile-spamassassin-parser-be && \
    mkdir -p /app && \
    cp ./bin/spamassassin-parser /app/spamassassin-parser

COPY ./build/docker/spamassassin-parser/entrypoint.sh /app/entrypoint.sh

FROM alpine:3.14.2 as deployment-container
LABEL maintainer="oleg.balunenko@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/obalunenko/spamassassin-parser"
LABEL stage="dev"

# Configure least privilege user
ARG UID=1000
ARG GID=1000
RUN addgroup -S spamassassin -g ${UID} && \
    adduser -S spamassassin -u ${GID} -G spamassassin -h /home/spamassassin -s /bin/sh -D spamassassin

WORKDIR /

ARG APK_CA_CERTIFICATES_VERSION=20191127-r5
RUN apk update && \
    apk add --no-cache \
        "ca-certificates=${APK_CA_CERTIFICATES_VERSION}" && \
    rm -rf /var/cache/apk/*

## Add the wait script to the image
ARG WAIT_VERSION=2.9.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/${WAIT_VERSION}/wait /wait
RUN chmod +x /wait

COPY --from=build-container /app/ /

RUN mkdir -p /data/input && \
    mkdir -p /data/result && \
    mkdir -p /data/archive && \
    chown -R spamassassin:spamassassin /data

ENTRYPOINT ["sh", "-c", "/wait && /entrypoint.sh"]

USER spamassassin
