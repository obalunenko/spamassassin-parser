ARG DOCKER_REPO
ARG DOCKER_GO_BASE_TAG=latest
ARG DOCKER_ALPINE_BASE_TAG=latest
ARG STAGE_ENV="prod"
# hadolint ignore=DL3007
FROM ${DOCKER_REPO}spamassassin-go-base:${DOCKER_GO_BASE_TAG} as build-container
LABEL maintainer="oleg.balunenko@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/obalunenko/spamassassin-parser"
LABEL stage=${STAGE_ENV}

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

FROM ${DOCKER_REPO}spamassassin-alpine-base:${DOCKER_ALPINE_BASE_TAG} as deployment-container
LABEL maintainer="oleg.balunenko@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/obalunenko/spamassassin-parser"
LABEL stage=${STAGE_ENV}

# Configure least privilege user
ARG UID=1000
ARG GID=1000
RUN addgroup -S spamassassin -g ${UID} && \
    adduser -S spamassassin -u ${GID} -G spamassassin -h /home/spamassassin -s /bin/sh -D spamassassin

WORKDIR /

COPY --from=build-container /app/ /

RUN mkdir -p /data/input && \
    mkdir -p /data/result && \
    mkdir -p /data/archive && \
    chown -R spamassassin:spamassassin /data

ENTRYPOINT ["sh", "-c", "/wait && /entrypoint.sh"]

USER spamassassin
