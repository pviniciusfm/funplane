SHELL := /bin/bash
BIN_DIR := $(GOPATH)/bin
RELEASE_DIR := ./release
DEP := $(BIN_DIR)/dep
BINARY_NAME := fanauth
PLATFORMS := linux
os = $(word 1, $@)
BINARIES := $(wildcard release/$(BINARY_NAME)-*)
artifactory_user_id ?= 
artifactory_password ?= 
artifactory_url ?= http://repo.frg.tech/artifactory/cloud/fanplane
CIRCLE_BRANCH ?= develop
fanplane_path := jaxf-github.fanatics.corp/cloud/fanplane
PKG_VERSION ?= 1.1
ifeq ($(CIRCLE_BRANCH), master)
	DIST_VERSION = $(PKG_VERSION)
else
	DIST_VERSION = $(CIRCLE_BRANCH)-$(PKG_VERSION)
endif
GEN_SECRET := ./pkg/auth/secrets.go
VENDOR_DIR := ./vendor
BASE_BINARY := $(RELEASE_DIR)/$(BINARY_NAME)

$(DEP):
	@go get -v -u github.com/golang/dep/cmd/dep

$(VENDOR_DIR): $(DEP)
	@dep ensure -v

dependencies: $(VENDOR_DIR) ## Installs dep and download external dependencies

.PHONY: clean
clean: ## Cleans project by removing releases folder ensuring a clean build and release
	rm -rf $(RELEASE_DIR)
	mkdir -p release

$(GEN_SECRET):
	go generate -v $(PKGS)

gen: $(GEN_SECRET) ## Generates dynamic go files (used only for secrets)

.PHONY: test
test: lint
	go test -v $(PKGS)

.PHONY: todo
todo: ## Show TODO task list present on the code
	@grep -IR --exclude-dir=vendor FUTURE: .| grep -v grep 
	@grep -IR --exclude-dir=vendor TODO: .| grep -v grep \
		&& echo "Solve all TODO comments before merge master!!!" \
		&& [ ${CIRCLE_BRANCH} = master ] && exit 1 \
		|| true

windows: CC=x86_64-w64-mingw32-gcc
darwin: CC=o64-clang
.PHONY: $(PLATFORMS)
$(PLATFORMS): clean dependencies gen
	LDFLAGS="${LDFLAGS} -linkmode external -s" GOOS=$(os) GOARCH=amd64 CC=$(CC) CGO_ENABLED=1 \
	go build -ldflags="-s -w -X $(fanplane_path)/cmd.version=$(PKG_VERSION) -X $(fanplane_path)/cmd.gitSha1=$(CIRCLE_SHA1) -X $(fanplane_path)/pkg/version.queryUser=$(aql_reader) -X $(fanplane_path)/pkg/version.queryPwd=$(aql_password)" -o $(BASE_BINARY)-$(os)-amd64
	@echo Running UPX in $(os) binary
	./ci/upx --brute --no-progress $(BASE_BINARY)-$(os)-amd64

$(BASE_BINARY): dependencies gen
	@echo "Building binary for the current OS"
	@mkdir -p $(RELEASE_DIR)
	@go build -ldflags="-s -w" -o $(BASE_BINARY)
	@echo "(_8(D) Done! Try it:"
	@echo "     $(BASE_BINARY) --help"

build: $(BASE_BINARY) ## Build for local environment only

.PHONY: rebuild
rebuild: clean build ## Force the build for local environment only

.PHONY: release
release: linux ## Build linux compatible binaries
	
.PHONY: docker-release
docker-release: ## Run release task inside a docker container for full cross os building
	@docker run --rm -w /go/src/jaxf-github.fanatics.corp/cloud/fanplane \
					 -v "$(shell pwd)":"/go/src/jaxf-github.fanatics.corp/cloud/fanplane" \
					 dockercore/golang-cross:1.11.0 \
					 make release

.PHONY: $(BINARIES)
$(BINARIES):
	@echo "Uploading $@ to artifactory."
	$(eval os := $(shell echo $@ | cut -d"-" -f2))
	@if [ $(os) == "windows" ]; \
	then curl -fSs  -u"$(artifactory_user_id):$(artifactory_password)" -T "$@" "$(artifactory_url)/$(DIST_VERSION)/$(os)/$(BINARY_NAME).exe;branch=$(CIRCLE_BRANCH);os=$(os);version=$(DIST_VERSION)"; \
	else curl -fSs  -u"$(artifactory_user_id):$(artifactory_password)" -T "$@" "$(artifactory_url)/$(DIST_VERSION)/$(os)/$(BINARY_NAME);branch=$(CIRCLE_BRANCH);os=$(os);version=$(DIST_VERSION)"; fi

.PHONY: publish
publish: $(BINARIES) ## Upload all released binaries files to artifactory

.PHONY: install
install: clean dependencies gen build ## Create the fanplane executable in $GOPATH/bin directory.
	install -m 0755 $(RELEASE_DIR)/$(BINARY_NAME) $(BIN_DIR)/$(BINARY_NAME)

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
