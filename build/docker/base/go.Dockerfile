FROM golang:1.20-alpine3.17
LABEL maintainer="oleg.balunenko@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/obalunenko/spamassassin-parser"
LABEL stage="base"

ARG PROJECT_URL=github.com/obalunenko/spamassassin-parser
RUN mkdir -p "${GOPATH}/src/${PROJECT_URL}/base-tools"

WORKDIR "${GOPATH}/src/${PROJECT_URL}/base-tools"

ARG APK_BASH_VERSION=~5
ARG APK_GIT_VERSION=~2
ARG APK_MAKE_VERSION=~4
ARG APK_OPENSSH_VERSION=~9
ARG APK_GCC_VERSION=~12
ARG APK_BUILDBASE_VERSION=~0
ARG APK_CA_CERTIFICATES_VERSION=20220614-r4
ARG APK_BINUTILS_VERSION=~2

RUN apk add --no-cache \
    "bash=${APK_BASH_VERSION}" \
	"git=${APK_GIT_VERSION}" \
	"make=${APK_MAKE_VERSION}" \
	"openssh-client=${APK_OPENSSH_VERSION}" \
	"build-base=${APK_BUILDBASE_VERSION}" \
    "gcc=${APK_GCC_VERSION}" \
    "ca-certificates=${APK_CA_CERTIFICATES_VERSION}" \
    "binutils-gold=${APK_BINUTILS_VERSION}"

ENV GOBIN="${GOPATH}/bin"
ENV PATH="${PATH}":"${GOBIN}"
