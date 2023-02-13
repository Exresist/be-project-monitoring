#!/bin/sh
FROM golang:latest as builder
RUN apk --no-cache add gcc g++ make git
WORKDIR $GOPATH/src/app/

COPY . .
RUN go mod init webserver
RUN go mod tidy

RUN export CGO_ENABLED=0 && make build

# Stage 2 - build final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/app/bin/be-project-monitoring go/bin/be-project-monitoring

# Run the binary.
EXPOSE 80
ENTRYPOINT /go/bin/be-project-monitoring --port 80
