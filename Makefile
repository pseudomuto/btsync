.DEFAULT_GOAL := build

BUILD_DEPS =
TEST_DEPS = $(BUILD_DEPS)

################################################################################
# PHONY Tasks
################################################################################
.PHONY: build retool setup test

build: $(BUILD_DEPS)
	$(info building binary bin/btsync...)
	@go build -o bin/btsync ./cmd/btsync

retool:
	@if test -z $(shell which retool); then \
		go get github.com/twitchtv/retool; \
		retool add github.com/golang/dep/cmd/dep v0.4.1; \
	fi

setup: retool
	$(info Synching dev tools and dependencies...)
	@retool sync
	@retool do dep ensure

test: $(TEST_DEPS)
	@go test -race -cover ./cmd/...
