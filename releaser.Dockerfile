FROM alpine
RUN apk add -U --no-cache ca-certificates


RUN mkdir -p /data/input && \
    mkdir -p /data/result && \
    mkdir -p /data/archive

COPY spamassassin-parser /

ENTRYPOINT ["/spamassassin-parser"]