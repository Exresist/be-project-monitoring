# Stage 1 - build executable in go container
FROM golang:latest AS builder

WORKDIR $GOPATH/src/app
COPY . .

RUN export CGO_ENABLED=0 && GOOS=linux go build -o /Users/HP/go/bin/app

FROM alpine:latest

ENTRYPOINT ["/Users/HP/go/bin/app"]
