FROM alpine:3.17.0
LABEL maintainer="oleg.balunenko@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/obalunenko/spamassassin-parser"

LABEL stage="base"

ARG APK_CA_CERTIFICATES_VERSION=~20230506
RUN apk update && \
    apk add --no-cache \
        "ca-certificates=${APK_CA_CERTIFICATES_VERSION}" && \
    rm -rf /var/cache/apk/*

WORKDIR /

## Add the wait script to the image
ARG WAIT_VERSION=2.9.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/${WAIT_VERSION}/wait /wait
RUN chmod +x /wait
