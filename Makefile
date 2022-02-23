## --------------------------------------
## Variables
## --------------------------------------

VUE_DIR := web
GO_BIN_OUT := main

## --------------------------------------
## All
## --------------------------------------

.PHONY: all
all: npm-install build

.PHONY: build
build: build-web build-go

.PHONY: clean
clean:
	rm -rf $(VUE_DIR)/build $(VUE_DIR)/node_modules $(GO_BIN_OUT) tmp

## --------------------------------------
## Vue and Node
## --------------------------------------

.PHONY: npm-install
npm-install: $(VUE_DIR)/package.json
	npm install --prefix ./$(VUE_DIR)

.PHONY: build-web
build-web:
	npm run --prefix ./$(VUE_DIR) build

.PHONY: npm-serve
npm-serve: $(VUE_DIR)/package.json
	npm install --prefix ./$(VUE_DIR)

## --------------------------------------
## Go
## --------------------------------------

.PHONY: build-go
build-go:
	go build -o $(GO_BIN_OUT)

.PHONY: run
run: build
	./$(GO_BIN_OUT)

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