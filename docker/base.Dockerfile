FROM golang:1.22.3-alpine3.19

RUN apk update && \
    apk add --no-cache \
    git \
    openssh \
    gcc \
    libc-dev \
    ca-certificates

WORKDIR /go/src/warehouse

RUN mkdir -p /root/.ssh \
    && ssh-keyscan github.com > /root/.ssh/known_hosts \
    && chmod 0700 -R /root/.ssh
