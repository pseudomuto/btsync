.DEFAULT_GOAL := build

BUILD_DEPS =
TEST_DEPS = $(BUILD_DEPS) pkg/internal/mocks/Connection.go

pkg/internal/mocks/Connection.go: pkg/bt/connection.go
	@retool do mockery -dir pkg/bt -name Connection --output pkg/internal/mocks

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

run:
	@docker run -p 9035:9035 shopify/bigtable-emulator

setup: retool
	$(info Synching dev tools and dependencies...)
	@retool sync
	@retool do dep ensure

test: $(TEST_DEPS)
	@go test -race -cover ./cmd/... ./pkg/...
