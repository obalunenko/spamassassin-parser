FROM golang:1.17.6-alpine3.15
LABEL maintainer="oleg.balunenko@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/obalunenko/spamassassin-parser"
LABEL stage="base"

ARG PROJECT_URL=github.com/obalunenko/spamassassin-parser
RUN mkdir -p "${GOPATH}/src/${PROJECT_URL}/base-tools"

WORKDIR "${GOPATH}/src/${PROJECT_URL}/base-tools"

ARG APK_GIT_VERSION=~2
ARG APK_NCURSES_VERSION=~6
ARG APK_MAKE_VERSION=~4
ARG APK_GCC_VERSION=~10
ARG APK_BASH_VERSION=~5.1
ARG APK_CURL_VERSION=~7
ARG APK_MUSL_DEV_VERSION=~1
ARG APK_UNZIP_VERSION=~6
ARG APK_CA_CERTIFICATES_VERSION=20211220-r0
ARG APK_LIBSTDC_VERSION=~10
ARG APK_BINUTILS_VERSION=~2
RUN apk update && \
    apk add --no-cache \
        "git=${APK_GIT_VERSION}" \
        "make=${APK_MAKE_VERSION}" \
        "gcc=${APK_GCC_VERSION}" \
        "bash=${APK_BASH_VERSION}" \
        "curl=${APK_CURL_VERSION}" \
        "musl-dev=${APK_MUSL_DEV_VERSION}" \
        "unzip=${APK_UNZIP_VERSION}" \
        "ca-certificates=${APK_CA_CERTIFICATES_VERSION}" \
        "libstdc++=${APK_LIBSTDC_VERSION}" \
        "binutils-gold=${APK_BINUTILS_VERSION}" && \
    rm -rf /var/cache/apk/*

# Get and install glibc for alpine
ARG APK_GLIBC_VERSION=2.29-r0
ARG APK_GLIBC_FILE="glibc-${APK_GLIBC_VERSION}.apk"
ARG APK_GLIBC_BIN_FILE="glibc-bin-${APK_GLIBC_VERSION}.apk"
ARG APK_GLIBC_BASE_URL="https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${APK_GLIBC_VERSION}"
# hadolint ignore=DL3018
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget -nv "${APK_GLIBC_BASE_URL}/${APK_GLIBC_FILE}" \
    && apk --no-cache add "${APK_GLIBC_FILE}" \
    && wget -nv "${APK_GLIBC_BASE_URL}/${APK_GLIBC_BIN_FILE}" \
    && apk --no-cache add "${APK_GLIBC_BIN_FILE}" \
    && rm glibc-*

COPY .git .git
COPY scripts scripts
COPY tools tools

COPY Makefile Makefile

# install tools from vendor
RUN make install-tools && \
    rm -rf "${GOPATH}/src/${PROJECT_URL}/base-tools"

ENV GOBIN="${GOPATH}/bin"
ENV PATH="${PATH}":"${GOBIN}"
