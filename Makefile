GOPATH:=$(shell go env GOPATH)
VERSION:=tmp-$(shell git describe --abbrev=8 --tags --always --dirty)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

ARCH :=	amd64
IMAGE_REGISTRY ?= docker-reg.github.com:5000/k8s/service
PROJECT_NAME := $(shell pwd | awk -F '/' '{print $$NF}')
ORG := xdev
IMAGE_NAME ?= $(ORG)-$(PROJECT_NAME)-$(ARCH)


ifeq ($(IMAGE_REGISTRY),)
IMAGE := $(IMAGE_NAME)
else
IMAGE := $(IMAGE_REGISTRY)/$(IMAGE_NAME)
endif


APP_NAME := $(ORG)-$(PROJECT_NAME)
HELM_RELEASE_NAME := $(APP_NAME)
HELM_CHART_NAME := github-chartmuseum/$(ORG)-service
HELM_CHART_DIR := $(shell cd ./deploy/helm-charts && pwd)
#DEPLOYMENT := $(ORG)
HELM_VALUES_FILE := $(HELM_CHART_DIR)/$(DEPLOYMENT)-values.yaml
HELM_RELEASE_NAMESPACE := $(ORG)

.PHONY: init

# init tool
init:
	GOPROXY=https://goproxy.cn,direct
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/nicksnyder/go-i18n/v2/goi18n@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/onsi/ginkgo/v2/ginkgo@latest

.PHONY: api
# generate api
api:
	$(shell protoc --proto_path=./api \
    	       --proto_path=./third_party \
     	       --go_out=paths=source_relative:./api \
     	       --go-grpc_out=paths=source_relative:./api \
     	       --go-errors_out=paths=source_relative:./api \
    	       $(API_PROTO_FILES))
	@for f in $(shell find api -name *.pb.go); do protoc-go-inject-tag -input=$${f} ;done

.PHONY: generate
# generate
generate:
	go mod tidy
	GOFLAGS=-mod=mod go generate ./...

.PHONY: wire
# generate DI code
wire:
	cd cmd/server && wire

.PHONY: i18n
# generate i18n
i18n:
	cd i18n && goi18n merge active.*.toml translate.*.toml && go run internal/i18n/generate.go

.PHONY: all
# generate all
all:
	make api;
	make generate;

.PHONY: lint
# lint
lint:
	@echo "===========> Run lint to lint source codes"
	golangci-lint run

.PHONY: fmt
# fmt
fmt:
	@echo "===========> Run fmt to check source codes"
	go fmt ./...

# build
.PHONY: build
build:
	@echo "===========> Run go build"
	mkdir -p bin && CGO_ENABLED=0 go build -o ./bin/ ./...
	chmod +x ./bin/server


Package ?= ./...
Method ?= .
.PHONY: test
test:
	@echo "===========> Run test to run Unit Test"
	@echo "DOCKER_MACHINE_NAME:$(DOCKER_MACHINE_NAME)"
	@echo "DOCKER_HOST:$(DOCKER_HOST)"
	@echo "DOCKER_URL:$(DOCKER_URL)"
	@echo "DOCKER_CERT_PATH:$(DOCKER_CERT_PATH)"
	ls /var/run/docker.sock
	go test --count=1 -v  $(Package) -test.run  $(Method) -covermode=set -coverprofile=report.txt

.PHONY:test-local
test-local:
	go test --count=1 -v  $(Package) -test.run  $(Method) -covermode=set -coverprofile=report.txt



# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
