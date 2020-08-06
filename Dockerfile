FROM golang:latest

RUN mkdir -p /go/src/github.com/hide168/charaboti.com

WORKDIR /go/src/github.com/hide168/charaboti.com

ADD . /go/src/github.com/hide168/charaboti.com

EXPOSE 8080