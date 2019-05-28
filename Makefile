BINARY=enc

.PHONY: all

all: build

build:
	go build -o ${BINARY} cmd/main.go