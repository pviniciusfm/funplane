 SHELL := /bin/bash
.DEFAULT_GOAL := build
export GO111MODULE := off
PKG_VERSION ?= dev
CIRCLE_SHA1 ?= $(shell git rev-parse HEAD)
TEST_RESULTS ?= ${PWD}/reports
BIN_DIR := $(GOPATH)/bin
DEP := $(BIN_DIR)/dep
VENDOR_DIR := ./vendor
GOMETALINTER := $(shell which golangci-lint) 
SOURCE_FILES ?= $(shell go list ./... | sed "s:^:$(GOPATH)/src/:")
OS := linux
fanplane_path := github.frg.tech/cloud/fanplane
LDFLAGS_REL := -ldflags="-s -w \
	-X $(fanplane_path)/cmd.version=$(PKG_VERSION) \
	-X $(fanplane_path)/cmd.gitSha1=$(CIRCLE_SHA1)"
PATH := ./bin:$(PATH)

$(DEP):
	@go get -v -u github.com/golang/dep/cmd/dep

$(VENDOR_DIR): $(DEP)
	@dep ensure -v

.PHONY: ci-test
ci-test: lint test ## Trigger tests ci cycle

.PHONY: setup
setup: $(DEP) ## Download tools required to check and build fanplane
	go get github.com/jstemmer/go-junit-report
	go get github.com/campoy/jsonenums
ifeq ($(OS,,), Darwin)
	brew install golangci-lint
else
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
endif

.PHONY: build
build: $(VENDOR_DIR) ## Download dependencies and compile project
	@go generate ./...
	@CGO_ENABLED=0 LDFLAGS="${LDFLAGS} -linkmode external -s" go build $(LDFLAGS_REL) -o fanplane main.go

.PHONY: lint
lint: ## Runs go lint checks
	@golangci-lint run ./...

.PHONY: e2e
e2e: build ## Compiles docker and use kubectl apply to deploy latest dev fanplane
	docker-compose build
	docker-compose push
	kubectl apply -f samples/kube

.PHONY: test
test: ## Run unit tests and output coverage report
	mkdir -p $(TEST_RESULTS)
	go test -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out ./... | go-junit-report -set-exit-code > test-results.xml
	go tool cover -html=coverage.out -o coverage.html
	mv test-results.xml coverage.out coverage.html ${TEST_RESULTS}

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
