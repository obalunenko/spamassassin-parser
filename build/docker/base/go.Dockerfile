FROM golang:1.17.3-alpine3.14
MAINTAINER oleg.balunenko@gmail.com

RUN mkdir -p ${GOPATH}/src/base-tools

WORKDIR ${GOPATH}/src/base-tools

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git make gcc bash curl musl-dev unzip ca-certificates libstdc++

RUN rm -rf /var/cache/apk/*

# Get and install glibc for alpine
ARG APK_GLIBC_VERSION=2.29-r0
ARG APK_GLIBC_FILE="glibc-${APK_GLIBC_VERSION}.apk"
ARG APK_GLIBC_BIN_FILE="glibc-bin-${APK_GLIBC_VERSION}.apk"
ARG APK_GLIBC_BASE_URL="https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${APK_GLIBC_VERSION}"
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget "${APK_GLIBC_BASE_URL}/${APK_GLIBC_FILE}" \
    && apk --no-cache add "${APK_GLIBC_FILE}" \
    && wget "${APK_GLIBC_BASE_URL}/${APK_GLIBC_BIN_FILE}" \
    && apk --no-cache add "${APK_GLIBC_BIN_FILE}" \
    && rm glibc-*

COPY .git .git
COPY scripts scripts
COPY tools tools

COPY Makefile Makefile

# install tools from vendor
RUN make install-tools

RUN rm -rf ${GOPATH}/src/base-tools

ENV GOBIN=${GOPATH}/bin
ENV PATH=${PATH}:${GOBIN}