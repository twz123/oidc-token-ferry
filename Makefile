PROGRAM       = oidc-token-ferry
GO_PACKAGE    = github.com/twz123/$(PROGRAM)

# binaries
GO  = go

VERSION              := $(shell git describe --tags --always)
GIT_UNTRACKEDCHANGES := $(shell git status --porcelain)
ifneq ($(GIT_UNTRACKEDCHANGES),)
	VERSION := $(VERSION)-dirty
endif

.PHONY: build
build: Makefile go.mod go.sum $(shell find pkg/ cmd/ -type f -name \*.go -print)
	$(GO) build -ldflags="-s -w -X $(GO_PACKAGE)/cmd/$(PROGRAM)/version.VERSION=$(VERSION)" ./cmd/$(PROGRAM)
