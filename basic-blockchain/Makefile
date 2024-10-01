BINARY := basic-blockchain

all: build test

build:deps
	@echo "====> Go build"
	@go build -o ${BINARY}

deps:
	@go mod tidy

test:
	./${BINARY} printchain
	./${BINARY} addblock -data "send 1 btc to van"
	./${BINARY} addblock -data "send 2 btc to alex"
	./${BINARY} addblock -data "send 3 more btc to van"
	./${BINARY} printchain

.PHONY: build deps test
