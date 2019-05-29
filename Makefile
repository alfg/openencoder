BINARY=enc
BINARY_SERVER=server
BINARY_WORKER=worker

.PHONY: all

all: build

build:
	go build -o ${BINARY} cmd/main.go

server:
	go build -o ${BINARY_SERVER} cmd/server/server.go

worker:
	go build -o ${BINARY_WORKER} cmd/worker/worker.go
