FROM golang:latest

RUN mkdir /go/src/charaboti.com

WORKDIR /go/src/charaboti.com

ADD . /go/src/charaboti.com

EXPOSE 8080