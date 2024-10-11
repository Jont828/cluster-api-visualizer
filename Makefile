## --------------------------------------
## Variables
## --------------------------------------

VUE_DIR := web
NODE_MODULES := ./$(VUE_DIR)/node_modules
DIST_FOLDER := ./$(VUE_DIR)/dist
GO_BIN_OUT := main

TAG ?= latest
REGISTRY ?= ghcr.io/jont828
IMAGE_NAME ?= cluster-api-visualizer
DOCKER_IMAGE ?= $(REGISTRY)/$(IMAGE_NAME)

ARCH ?= $(shell go env GOARCH)
# ALL_ARCH = amd64 arm64
ALL_ARCH = amd64 arm arm64 ppc64le s390x

# Build time versioning details.
LDFLAGS := $(shell hack/version.sh)

## --------------------------------------
## All
## --------------------------------------

# Default target is to build and run the app
.PHONY: all
all: npm-install build run

.PHONY: build
build: build-web build-go

.PHONY: clean
clean:
	rm -rf $(DIST_FOLDER) $(NODE_MODULES) $(GO_BIN_OUT) tmp

## --------------------------------------
## Vue and Node
## --------------------------------------

.PHONY: npm-install
npm-install: $(VUE_DIR)/package.json
	npm install --prefix ./$(VUE_DIR)

.PHONY: build-web
build-web: $(NODE_MODULES)
	npm run --prefix ./$(VUE_DIR) build

.PHONY: npm-serve
serve: $(VUE_DIR)/package.json $(NODE_MODULES)
	npm run --prefix ./$(VUE_DIR) serve

.PHONY: clean-dist
clean-dist:
	rm -rf $(DIST_FOLDER)

## --------------------------------------
## Go
## --------------------------------------

.PHONY: build-go
build-go:
	go build -ldflags "$(LDFLAGS)" -o $(GO_BIN_OUT)

.PHONY: run
run: $(GO_BIN_OUT) $(DIST_FOLDER)
	./$(GO_BIN_OUT)

.PHONY: go-run
go-run: $(DIST_FOLDER)
	go run -ldflags "$(LDFLAGS)" main.go 

.PHONY: air
air: .air.toml
	air

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy

.PHONY: go-vet
go-vet:
	go vet ./...

.PHONY: go-fmt
go-fmt:
	go fmt ./...

## --------------------------------------
## Docker
## --------------------------------------

.PHONY: build-and-deploy docker-build
build-and-deploy:
	docker image tag ghcr.io/jont828/cluster-api-visualizer-$(ARCH)\:v1.3.1 ghcr.io/jont828/cluster-api-visualizer\:v1.3.1
	kind load docker-image -n hmc-dev ghcr.io/jont828/cluster-api-visualizer\:v1.3.1
	./hack/deploy-local-to-kind.sh

.PHONY: docker-build-all
docker-build-all: $(addprefix docker-build-,$(ALL_ARCH)) 

docker-build-%:
	$(MAKE) ARCH=$* docker-build

.PHONY: docker-build
docker-build: 
	docker build --no-cache --build-arg ARCH=$(ARCH) --build-arg ldflags="$(LDFLAGS)" -t $(DOCKER_IMAGE)-$(ARCH):$(TAG) .

.PHONY: docker-push
docker-push: 
	docker push $(DOCKER_IMAGE)-$(ARCH):$(TAG)

.PHONY: docker-push-all
docker-push-all: $(addprefix docker-push-,$(ALL_ARCH)) ## Push all the architecture docker images.
	$(MAKE) docker-push-manifest

docker-push-%:
	$(MAKE) ARCH=$* docker-push

.PHONY: docker-push-manifest
docker-push-manifest: ## Push the fat manifest docker image.
	## Minimum docker version 18.06.0 is required for creating and pushing manifest images.
	docker manifest create --amend $(DOCKER_IMAGE):$(TAG) $(shell echo $(ALL_ARCH) | sed -e "s~[^ ]*~$(DOCKER_IMAGE)\-&:$(TAG)~g")
	@for arch in $(ALL_ARCH); do docker manifest annotate --arch $${arch} ${DOCKER_IMAGE}:${TAG} ${DOCKER_IMAGE}-$${arch}:${TAG}; done
	docker manifest push --purge ${DOCKER_IMAGE}:${TAG}

## --------------------------------------
## Helm
## --------------------------------------

.PHONY: update-helm
update-helm:
	./hack/update-helm-repo.sh