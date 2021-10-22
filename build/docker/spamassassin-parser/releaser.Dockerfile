FROM alpine:3.14.2
# Configure least privilege user
ARG UID=1000
ARG GID=1000
RUN addgroup -S spamassassin -g ${UID} && \
    adduser -S spamassassin -u ${GID} -G spamassassin -h /home/spamassassin -s /bin/sh -D spamassassin

WORKDIR /

ARG APK_CA_CERTIFICATES_VERSION=20191127-r5
RUN apk update && \
    apk add --no-cache \
        "ca-certificates=${APK_CA_CERTIFICATES_VERSION}" \
    rm -rf /var/cache/apk/*

RUN mkdir -p /data/input && \
    mkdir -p /data/result && \
    mkdir -p /data/archive && \
    chown -R spamassassin:spamassassin /data

COPY spamassassin-parser /
COPY build/docker/spamassassin-parser/entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]

USER spamassassin