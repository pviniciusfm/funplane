FROM golang:1.11.0-alpine3.8

COPY ./fanplane /go/bin
WORKDIR /go

ENTRYPOINT ["/go/bin/fanplane"]
