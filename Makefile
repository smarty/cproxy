#!/usr/bin/make -f

test: fmt
	GORACE="atexit_sleep_ms=50" go test -timeout=1s -race -covermode=atomic -count=1 ./...

fmt:
	go fmt ./...

compile:
	go build ./...

build: test compile

.PHONY: test fmt compile build
