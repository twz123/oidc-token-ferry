PROGRAM       = oidc-token-ferry
GO_PACKAGE    = github.com/twz123/$(PROGRAM)
BUILDER_IMAGE = docker.io/golang:1.9.2-alpine3.7

# binaries
DOCKER = docker
GO     = go
DEP    = dep

OS_ARCH_PROGRAMS =
PROGRAM_DEPENDENCIES = Makefile Gopkg.lock $(shell find pkg/ cmd/ -type f -name \*.go -print)

$(PROGRAM): $(PROGRAM_DEPENDENCIES)
	$(GO) build ./cmd/oidc-token-ferry

define _os_arch_program =
OS_ARCH_PROGRAMS += $(PROGRAM).$(1)-$(2)
oidc-token-ferry.$(1)-$(2): $(PROGRAM_DEPENDENCIES)
	$(DOCKER) run --rm -e GOOS=$(1) -e GOARCH=$(2) -e CGO_ENABLED=0 -v "$(shell pwd -P):/go/src/$(GO_PACKAGE):ro" -w "/go/src/$(GO_PACKAGE)/cmd/$(PROGRAM)" $(BUILDER_IMAGE) \
	sh -c 'go build -ldflags="-s -w" -o /tmp/go.out && apk add --no-cache upx 1>&2 && upx -o /tmp/go.out.upx /tmp/go.out 1>&2 && cat /tmp/go.out.upx' > $(PROGRAM).$(1)-$(2) || { rm $(PROGRAM).$(1)-$(2); exit 1; }
	chmod +x $(PROGRAM).$(1)-$(2)
endef

$(eval $(call _os_arch_program,linux,amd64))
$(eval $(call _os_arch_program,darwin,amd64))

.PHONY: all
all: $(OS_ARCH_PROGRAMS)

Gopkg.lock: Gopkg.toml $(shell find vendor/ -type f -name \*.go -print)
	$(DEP) ensure
	touch Gopkg.lock

.PHONY: clean
clean:
	rm -f $(PROGRAM) $(OS_ARCH_PROGRAMS)
