CC=go

all: build

install:
	$(CC) get

build: main.go
	$(CC) build ./*.go

run:
	$(CC) run ./**.go

run-build: build
	./sync-pubsub

clean-build: install build
