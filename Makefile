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

##@ All:

# Default target is to build and run the app
.PHONY: all
all: npm-install build run ## Install the npm dependencies, build the Vue app, build the Go binary and run the app.

.PHONY: build
build: build-web build-go ## Build the Vue app and the Go binary.

.PHONY: clean
clean: ## Remove the dist folder, node_modules, the Go binary and the tmp folder.
	rm -rf $(DIST_FOLDER) $(NODE_MODULES) $(GO_BIN_OUT) tmp

## --------------------------------------
## JavaScript
## --------------------------------------

##@ JavaScript:

.PHONY: npm-install
npm-install: $(VUE_DIR)/package.json ## Install the npm dependencies.
	npm install --prefix ./$(VUE_DIR)

.PHONY: build-web
build-web: $(NODE_MODULES) ## Build the Vue app.
	npm run --prefix ./$(VUE_DIR) build

.PHONY: npm-serve
serve: $(VUE_DIR)/package.json $(NODE_MODULES) ## Run the Vue app.
	npm run --prefix ./$(VUE_DIR) serve

.PHONY: clean-dist
clean-dist: ## Remove the dist folder.
	rm -rf $(DIST_FOLDER)

## --------------------------------------
## Go
## --------------------------------------

##@ Go:

.PHONY: build-go
build-go: ## Build the Go binary.
	go build -ldflags "$(LDFLAGS)" -o $(GO_BIN_OUT)

.PHONY: run
run: $(GO_BIN_OUT) $(DIST_FOLDER) ## Run the Go app using a binary.
	./$(GO_BIN_OUT)

.PHONY: go-run
go-run: $(DIST_FOLDER) ## Run the Go app.
	go run -ldflags "$(LDFLAGS)" main.go 

.PHONY: air
air: .air.toml ## Start `air`.
	air

.PHONY: go-mod-tidy
go-mod-tidy: ## Run `go mod tidy`.
	go mod tidy

.PHONY: go-vet
go-vet: ## Run `go vet`.
	go vet ./...

.PHONY: go-fmt
go-fmt: ## Run `go fmt`.
	go fmt ./...

## --------------------------------------
## Docker
## --------------------------------------

.PHONY: docker-build-all
docker-build-all: $(addprefix docker-build-,$(ALL_ARCH)) ## Build all the architecture docker images.

docker-build-%: ## Build the docker image for the specified architecture.
	$(MAKE) ARCH=$* docker-build

.PHONY: docker-build
docker-build: ## Build the docker image for the current architecture.
	docker build --no-cache --build-arg ARCH=$(ARCH) --build-arg ldflags="$(LDFLAGS)" -t $(DOCKER_IMAGE)-$(ARCH):$(TAG) .

.PHONY: docker-push
docker-push: ## Push the docker image for the current architecture.
	docker push $(DOCKER_IMAGE)-$(ARCH):$(TAG)

.PHONY: docker-push-all
docker-push-all: $(addprefix docker-push-,$(ALL_ARCH)) ## Push all the architecture docker images.
	$(MAKE) docker-push-manifest

docker-push-%: ## Push the docker image for the specified architecture.
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

##@ Helm:

.PHONY: update-helm
update-helm: ## Update the helm repo.
	./hack/update-helm-repo.sh

## --------------------------------------
## Help
## --------------------------------------

##@ Help:

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)