FROM golang:1.15-alpine AS builder
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app/
RUN apk add --no-cache git

ENV GO111MODULE="on"
ENV CGO_ENABLED=0

RUN go build -o="goapp"

FROM alpine:latest
RUN mkdir -p /home/app
WORKDIR /home/app
COPY --from=builder /go/src/app /home/app
ENTRYPOINT /home/app/goapp
