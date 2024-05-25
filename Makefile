GOPATH:=$(shell go env GOPATH)
GIT_KEY_PATH?=~/.ssh/id_rsa

.PHONY: init
init:
	@go install github.com/google/wire/cmd/wire@latest

build_base:
	docker build -t warehouse-base \
		--ssh default=${SSH_AUTH_SOCK} \
		-f docker/base.Dockerfile .

build_dev:
	docker build -t warehouse-dev \
		--ssh default=${SSH_AUTH_SOCK} \
		-f docker/dev.Dockerfile .

build_prod:
	docker build -t warehouse-production \
		--ssh default=${SSH_AUTH_SOCK} \
		-f docker/production.Dockerfile .

.PHONY: build
build:
	@go build -o warehouse cmd/main.go
