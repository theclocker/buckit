FROM golang:1.10

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN go install -v ./...

RUN go get github.com/codegangsta/gin
RUN go get -u github.com/theclocker/api_manager/...