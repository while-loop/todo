# service specific vars
SERVICE     := todo
VERSION     := 0.0.1

ORG         := toyotasupra
TARGET      := ${SERVICE}d
COMMIT      := $(shell git rev-parse --short HEAD)
BUILD_TIME  := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
IMAGE_NAME  := ${ORG}/${SERVICE}
PACKAGE 	:= $(shell pwd | sed "s,${GOPATH}/src/,,")

.PHONY: proto deps test build cont cont-nc all deploy help clean lint
.DEFAULT_GOAL := help

help: ## halp
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: deps lint test build ## get && test && build

build: clean lint ## build service binary file
	@echo "[build] building go binary"
	@go build \
		-ldflags "-s -w \
		-X ${PACKAGE}/pkg.Version=${VERSION} \
		-X ${PACKAGE}/pkg.Commit=${COMMIT} \
		-X ${PACKAGE}/pkg.BuildTime=${BUILD_TIME}" \
		-o ${GOPATH}/bin/${TARGET} ./cmd/${TARGET}
	@todod -v

clean: ## remove service bin from $GOPATH/bin
	@echo "[clean] removing service files"
	rm -f ${GOPATH}/bin/${TARGET}

cont: ## build a non-cached service container
	docker build -t ${IMAGE_NAME} -t ${IMAGE_NAME}:${VERSION} . --no-cache

cont-c: ## build a cached service container
	docker build -t ${IMAGE_NAME} -t ${IMAGE_NAME}:${VERSION} .

deploy: ## deploy lastest built container to docker hub
	docker push ${IMAGE_NAME}

deps: ## get service pkg + test deps
	@echo "[deps] getting go deps"
	go get -v -t ./...

lint: ## apply golint
	@echo "[lint] applying go fmt & vet"
	go fmt ./...
	go vet ./...

test: lint ## test service code
	@echo "[test] running tests w/ cover"
	go test ./... -cover
