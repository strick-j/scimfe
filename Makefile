GO ?= go

# Stub for plain 'make'
all:

include scripts/migrate.mk
include scripts/golangci.mk

.PHONY: run
run:
	go run ./cmd/scimfe -config /config/config.dev.yaml