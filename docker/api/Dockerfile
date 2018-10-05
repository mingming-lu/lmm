FROM golang:1.11

ENV GO111MODULE on

COPY . /go/src/lmm/api

WORKDIR /go/src/lmm/api

RUN go mod download

RUN go install lmm/api
