export GO111MODULE=on

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
OUT := bin/backoffice-cebuana-api
OUT_PATH=$(shell pwd)/bin/$(GOOS)_$(GOARCH)
TEST_PACKAGE=./internal/...
TEST_CONFIG=-race -timeout 600s -count=1 -p 1


clean:
	rm -rf ./bin/*
.PHONY: clean

clean.bin:
	rm -rf $(OUT_PATH)
.PHONY: clean.bin

build: clean
	go build -ldflags "-w -s" -o $(OUT) ./cmd/be-project-monitoring
.PHONY: build