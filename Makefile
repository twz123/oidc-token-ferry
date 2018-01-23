BUILD_DIR=/go/src/github.com/twz123/oidc-token-ferry
BUILDER_IMAGE=docker.io/golang:1.9.2-alpine3.7

DOCKER_IMAGE_NAME=quay.io/twz123/oidc-token-ferry
DOCKER_IMAGE_TAG=$(shell git describe --tags --always --dirty)

# binaries
DOCKER=docker
DEP=dep

oidc-token-ferry: Makefile Gopkg.lock $(shell find pkg/ cmd/ -type f -name \*.go -print)
	$(DOCKER) run --rm -e CGO_ENABLED=0 -v "$(shell pwd -P):$(BUILD_DIR):ro" -w "$(BUILD_DIR)/cmd/oidc-token-ferry" $(BUILDER_IMAGE) \
	go build -o /dev/stdout > oidc-token-ferry || { rm oidc-token-ferry; exit 1; }
	chmod +x oidc-token-ferry

Gopkg.lock: Gopkg.toml $(shell find vendor/ -type f -name \*.go -print)
	$(DEP) ensure
	touch Gopkg.lock

clean:
	rm -f oidc-token-ferry

.PHONY: dockerize
dockerize: oidc-token-ferry Dockerfile
	$(DOCKER) build . -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: publish-docker-image
publish-docker-image: dockerize
	$(DOCKER) push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
