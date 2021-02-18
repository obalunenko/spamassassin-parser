ARG DOCKER_REPO=olegbalunenko
FROM ${DOCKER_REPO}/spamassassin-go-base:latest as build-container

MAINTAINER Oleg Balunenko

ENV PROJECT_DIR=${GOPATH}/src/github.com/obalunenko/spamassassin-parser

RUN mkdir -p ${PROJECT_DIR}

WORKDIR ${PROJECT_DIR}

COPY .git .git
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY scripts scripts
COPY vendor vendor
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile

# vet project
RUN make vet
# compile executable
RUN make compile

RUN mkdir /app
RUN cp ./bin/spamassassin-parser /app/spamassassin-parser


FROM alpine:3.13 as deployment-container
RUN apk add -U --no-cache ca-certificates


RUN mkdir -p /data/input && \
    mkdir -p /data/result && \
    mkdir -p /data/archive

COPY --from=build-container /app/spamassassin-parser /spamassassin-parser

ENTRYPOINT ["/spamassassin-parser"]

