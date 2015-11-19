FROM docker.yb0t.cc/alpine-golang:1.5.1

MAINTAINER Yieldbot Infrastructure <infra@yieldbot.com>


RUN echo "http://dl-4.alpinelinux.org/alpine/v2.6/main" >> /etc/apk/repositories && \
    apk-install make=3.82-r6 bash man mdocml-apropos mdocml && \
    go get -u github.com/golang/lint/golint
