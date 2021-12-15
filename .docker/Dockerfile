FROM golang:1.17-alpine

WORKDIR /app

ADD ./service/go.mod ./
RUN go mod download


ADD ./main/*.go ./
RUN go get

RUN go get github.com/codegangsta/gin