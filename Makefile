ifeq ($(origin GITROOT), undefined)
GITROOT := $(shell git rev-parse --show-toplevel)
endif

# ENV
ENV     ?= dev

# Commands
DOCKER  ?= docker
GREP    ?= grep
MKDIR_P ?= mkdir -p
ABIGEN  ?= $(GITROOT)/bin/abigen

# GO Commands
GO            ?= go
GOIMPORTS     ?= goimports
MOCKGEN       ?= $(wildcard $(GITROOT)/bin/mockgen)
GOLANGCI_LINT ?= $(wildcard $(GITROOT)/bin/golangci-lint)
GOLINT        ?= $(GOLANGCI_LINT) run
PROTOC        ?= $(or $(wildcard $(GITROOT)/bin/protoc), $(shell which protoc))

GIT_COMMIT_HASH ?= $(shell git rev-parse HEAD | cut -c -10)
NOW ?= $(shell date +%Y%m%d%H%M%S)

# GO build flags
GO_IMPORT_PATH = github.com/popodidi/lambda-protocol
GO_BUILD_LDFLAGS =
GO_BUILD_LDFLAGS += -X $(GO_IMPORT_PATH)/pkg/version.version=$(GIT_COMMIT_HASH)
GO_BUILD_LDFLAGS += -X $(GO_IMPORT_PATH)/pkg/version.buildtime=$(NOW)

GO_BUILD_FULL_OPTS = $(GO_BUILD_OPTS) -ldflags "$(GO_BUILD_LDFLAGS)"

# Path
CMDDIR   ?= $(GITROOT)/cmd
BUILDDIR ?= $(GITROOT)/_build

GO_DEPS = \
	goimports;golang.org/x/tools/cmd/goimports \
	protoc-gen-go;github.com/golang/protobuf/protoc-gen-go \
	protoc-gen-gotag;github.com/srikrsna/protoc-gen-gotag \
	mockgen;github.com/golang/mock/mockgen

define INSTALL_RULE
install-$1:
ifeq (,$(shell which $1))
	$2
else
	@echo "$1 is installed"
endif
endef

define GO_INSTALL_RULE
$(eval $(call INSTALL_RULE,$1,GOBIN=$(GITROOT)/bin $(GO) get $2))
endef

define BUILD_RULE
$1: pre-build
	$(GO) build $(GO_BUILD_FULL_OPTS) -o $(BUILDDIR)/$1 $(CMDDIR)/$1
endef

define PROTOC_RULE
proto-$1:
	PATH="$(GITROOT)/bin:$(PATH)" $(PROTOC) -I=$(GITROOT) \
		--go_out=plugins=grpc:$(GITROOT) \
		$(GITROOT)/$1/*.proto

	PATH="$(GITROOT)/bin:$(PATH)" $(PROTOC) -I=$(GITROOT) \
		--gotag_out=auto="bson-as-snake":. \
		$(GITROOT)/$1/*.proto
endef

APPS = $(shell ls $(CMDDIR) | grep -v playground)

PROTO_DIRS =
# PROTO_DIRS = \
# 	pkg/engine

.PHONY: default apps pre-build lint test tidy dep install-golangci-lint clean $(APPS)

default: apps

# Dependencies
dep: install-golangci-lint $(foreach DEP,$(GO_DEPS), $(eval CMD = $(word 1,$(subst ;, ,$(DEP)))) install-$(CMD))
$(foreach DEP,$(GO_DEPS), \
	$(eval CMD = $(word 1,$(subst ;, ,$(DEP)))) \
	$(eval SRC = $(word 2,$(subst ;, ,$(DEP)))) \
	$(eval $(call GO_INSTALL_RULE,$(CMD),$(SRC))) \
)

install-golangci-lint:
ifeq (, $(GOLANGCI_LINT))
	@wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.23.8
else
	@echo "golangci-lint is installed"
endif

proto: $(foreach proto_dir, $(PROTO_DIRS), proto-$(proto_dir))
$(foreach proto_dir, $(PROTO_DIRS), $(eval $(call PROTOC_RULE,$(proto_dir))))

# Lint, and test
precommit: tidy apps format lint test

format:
	find . -name "*.go" | xargs $(GOIMPORTS) -w -local $(GO_IMPORT_PATH)

lint:
	$(GOLINT) --skip-dirs="(^|/)playground($|/)" $(GITROOT)/...

test: generate
	$(GO) test $(GITROOT)/...

test-ci:
	$(GO) test -tags=gitlabci $(GITROOT)/...

tidy:
	$(GO) mod tidy

# Build apps
pre-build:
	@$(MKDIR_P) $(BUILDDIR)

generate:
	PATH="$(GITROOT)/bin:$(PATH)" $(GO) generate ./...

clean:
	rm -rvf $(BUILDDIR)/*

apps: $(APPS)
$(foreach app, $(APPS), $(eval $(call BUILD_RULE,$(app))))
