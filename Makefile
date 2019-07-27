BINARY=openencoder

.PHONY: all

all: build

build:
    CGO_ENABLED=0 GOOS=linux go build -installsuffix 'static' -v -o ${BINARY} .