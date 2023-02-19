# Stage 1 - build executable in go container
FROM golang:latest as builder

WORKDIR $GOPATH/src/be-project-monitoring/
COPY . .

RUN export CGO_ENABLED=0 && make build

# Stage 2 - build final image
FROM alpine:latest

# Copy our static executable
COPY --from=builder /go/src/be-project-monitoring/bin/be-project-monitoring go/bin/be-project-monitoring

# Run the binary.
ENTRYPOINT ["go/bin/be-project-monitoring"]
