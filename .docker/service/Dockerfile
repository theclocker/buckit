FROM golang:1.17-alpine

RUN apk add --no-cache bash

WORKDIR /service
# ENV GOPATH=/service

# ADD ./service/go.mod ./
# RUN go mod download

ADD ./service/*.go ./
# RUN go get

# RUN go get github.com/codegangsta/gin