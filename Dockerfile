#!/bin/sh
FROM golang:latest as builder

WORKDIR $GOPATH/src/app/

COPY . .

RUN export CGO_ENABLED=0 && make build

# Stage 2 - build final image
FROM alpine:latest

COPY --from=builder /go/src/app/bin/be-project-monitoring go/bin/be-project-monitoring

# Run the binary.
ENTRYPOINT ["go/bin/be-project-monitoring"]
