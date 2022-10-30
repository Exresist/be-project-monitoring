GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
OUT := bin/be-project-monitoring
OUT_PATH=$(shell pwd)/bin/$(GOOS)_$(GOARCH)

clean:
	rm -rf ./bin/*
.PHONY: clean

clean.bin: ## remove $(OUT_PATH) directory
	rm -rf $(OUT_PATH)
.PHONY: clean.bin

build: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o $(OUT) ./cmd/be-project-monitoring
.PHONY: build

gen:
	go generate ./...
.PHONY: gen

dbuild:
	docker-compose -f ../be-project-monitoring/docker-compose.yml build
.PHONY: dbuild

dstart:
	docker-compose -f ../be-project-monitoring/docker-compose.yml up
.PHONY: dstart

run: dbuild dstart
.PHONY: run

lint: ## run linters for project
	$(OUT_PATH)/golangci-lint run
.PHONY: lint

prepare: clean install.tools ## performs steps needed before first build
.PHONY: prepare