VERSION ?= 2.0.0-dev

IMAGE_NAME := flexconstructor/openvidu-tutorial

PROJ_NAME := github.com/flexconstructor/openvidu-tutorial

BUILD_DIR := _build
RELEASE_DIR := _release

RELEASE_BRANCH := master
MAINLINE_BRANCH := dev
CURRENT_BRANCH := $(shell git branch | grep \* | cut -d ' ' -f2)

GOLANG_VER ?= 1.9

no-cache ?= no
port ?= !default
bench ?= !default
docker ?= yes


comma := ,
empty :=
space := $(empty) $(empty)
eq = $(if $(or $(1),$(2)),$(and $(findstring $(1),$(2)),\
                                $(findstring $(2),$(1))),1)



# Run project in Dockerized development environment.
#
# Usage:
#	make run

run: | build
	docker-compose down
	docker-compose up --build



# Resolve all project dependencies.
#
# Usage:
#	make deps

deps: | deps.glide deps.tools




# Resolve Golang Glide project dependencies.
#
# Usage:
#	make deps.glide [cmd=]

cmd ?= install

deps.glide:
	docker run --rm -v "$(PWD)":/go/src/$(PROJ_NAME) -w /go/src/$(PROJ_NAME) \
	                -v "$(PWD)/_cache/glide:/root/.glide/cache" \
		instrumentisto/glide $(cmd)




# Resolve Golang binary toolchain of project development dependencies.
#
# Usage:
#	make deps.tools

tools := \
	github.com/smartystreets/goconvey \
	github.com/alecthomas/gometalinter \
	github.com/yookoala/realpath \
	github.com/go-playground/overalls

deps.tools:
	mkdir -p vendor/_bin
	(set -e ; $(foreach tool, $(tools), \
		docker run --rm -v "$(PWD)/vendor":/go/src -w /go/src/$(tool) \
			golang:$(GOLANG_VER) \
				go build -o /go/src/_bin/$(word 3,$(subst /, ,$(tool))); \
	))
	docker run --rm -v "$(PWD)/vendor":/go/src -v "$(PWD)/vendor/_bin":/go/bin \
		golang:$(GOLANG_VER) \
			gometalinter --install




# Build project.
#
# Usage:
#	make build [VERSION=]

build:
	mkdir -p $(BUILD_DIR)
	rm -rf $(BUILD_DIR)/*
	docker run --rm -v "$(PWD)":/go/src/$(PROJ_NAME) -w /go/src/$(PROJ_NAME) \
		golang:$(GOLANG_VER) \
			go build -installsuffix netgo -tags='netgo jsoniter' \
			         -o $(BUILD_DIR)/openvidu_tutorial main.go
	printf "$(VERSION)" > $(BUILD_DIR)/version




# Format all project Golang sources with Gofmt.
#
# Usage:
#	make fmt

fmt:
	docker run --rm -v "$(PWD)":/go/src/$(PROJ_NAME) -w /go/src/$(PROJ_NAME) \
		golang:$(GOLANG_VER) \
			go fmt ./...



# Lint all project Golang sources with Go Meta Linter.
#
# Commands `go install .` and `go test -i` required for some Go Meta Linter
# tools to run correctly. Details may be found here:
# https://github.com/alecthomas/gometalinter/issues/91
#
# Usage:
#	make lint

lint:
ifneq ($(wildcard vendor/$(PROJ_NAME)),)
	rm -rf vendor/$(PROJ_NAME)
endif
	docker run --rm -v "$(PWD)":/go/src/$(PROJ_NAME) -w /go/src/$(PROJ_NAME) \
	                -v "$(PWD)/vendor":/go/src \
	                -v "$(PWD)/vendor/_bin":/go/bin \
		golang:$(GOLANG_VER) \
			gometalinter --config=.gometalinter.json ./...
ifneq ($(wildcard vendor/$(PROJ_NAME)),)
	rm -rf vendor/$(PROJ_NAME)
endif



# Run all project tests.
#
# Optional 'bench' parameter may be used to run Go benchmarks.
# It assumes the same values as `-bench` flag of `go test`.
# For example: `make test bench=.`.
#
# Usage:
#	make test [bench=]

test-bench-arg = $(if $(call eq,$(bench),!default),,-benchmem -bench=$(bench))

test:
	docker run --rm -v "$(PWD)":/go/src/$(PROJ_NAME) -w /go/src/$(PROJ_NAME) \
	                -v "$(PWD)/vendor":/go/src \
	                -v "$(PWD)/vendor/_bin":/go/bin \
		golang:$(GOLANG_VER) \
			overalls \
				-project=$(PROJ_NAME) \
				-covermode=atomic \
				-ignore='.git,vendor_tools,vendor,node_modules,_cache' \
				-- -race $(test-bench-arg)
	docker run --rm -v "$(PWD)":/go/src/$(PROJ_NAME) -w /go/src/$(PROJ_NAME) \
		golang:$(GOLANG_VER) \
			go tool cover -func=overalls.coverprofile



# Run GoConvey Web UI for project Golang tests.
#
# Usage:
#	make run.goconvey [port=<8080>] [docker=(yes|no)]

goconvey-port = $(if $(call eq,$(port),!default),8080,$(port))

run.goconvey:
ifeq ($(docker),yes)
	docker run --rm -v "$(PWD)":/go/src/$(PROJ_NAME) -w /go/src/$(PROJ_NAME) \
	                -v "$(PWD)/vendor/github.com/smartystreets/goconvey":/go/src/github.com/smartystreets/goconvey \
	                -v "$(PWD)/vendor/_bin/goconvey":/go/bin/goconvey \
	                -p $(goconvey-port):$(goconvey-port) \
		golang:$(GOLANG_VER) \
			goconvey \
				-host 0.0.0.0 -port $(goconvey-port) \
				-cover \
				-excludedDirs=".git,vendor_tools,vendor,node_modules,_cache"
else
	goconvey \
		-port $(goconvey-port) \
		-cover \
		-excludedDirs=".git,vendor_tools,vendor,node_modules,_cache"
endif



# Squash changes of the current Git branch onto another Git branch.
#
# WARNING: You must merge `onto` branch in the current branch before squash!
#
# Usage:
#	make squash [onto=] [del=(no|yes)]

onto ?= $(MAINLINE_BRANCH)
del ?= no
upstream ?= origin

squash:
ifeq ($(CURRENT_BRANCH),$(onto))
	@echo "--> Current branch is '$(onto)' already" && false
endif
	git checkout $(onto)
	git branch -m $(CURRENT_BRANCH) orig-$(CURRENT_BRANCH)
	git checkout -b $(CURRENT_BRANCH)
	git branch --set-upstream-to $(upstream)/$(CURRENT_BRANCH)
	git merge --squash orig-$(CURRENT_BRANCH)
ifeq ($(del),yes)
	git branch -d orig-$(CURRENT_BRANCH)
endif



.PHONY: deps deps.glide deps.tools build fmt test run.goconvey run
