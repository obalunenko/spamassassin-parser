ARG DOCKER_REPO
ARG DOCKER_ALPINE_BASE_TAG=latest

FROM ${DOCKER_REPO}spamassassin-alpine-base:${DOCKER_ALPINE_BASE_TAG} as deployment-container
LABEL maintainer="oleg.balunenko@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/obalunenko/spamassassin-parser"
LABEL stage="release"

# Configure least privilege user
ARG UID=1000
ARG GID=1000
RUN addgroup -S spamassassin -g ${UID} && \
    adduser -S spamassassin -u ${GID} -G spamassassin -h /home/spamassassin -s /bin/sh -D spamassassin

WORKDIR /

COPY spamassassin-parser /
COPY build/docker/spamassassin-parser/entrypoint.sh /

RUN mkdir -p /data/input && \
    mkdir -p /data/result && \
    mkdir -p /data/archive && \
    chown -R spamassassin:spamassassin /data

ENTRYPOINT ["/entrypoint.sh"]

USER spamassassin