GO ?= go

# Stub for plain 'make'
all:

include scripts/migrate.mk
include scripts/golangci.mk

.PHONY: run
run:
	go run ./cmd/scimfe -config /configs/config.dev.yaml

	.PHONY: e2e
e2e:
	go test -v -count=1 ./test/e2e/...