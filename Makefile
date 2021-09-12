CC=go

all: build

install:
	$(CC) get

build: main.go
	$(CC) build ./*.go

run:
	$(CC) run ./**.go

run-build: build
	./main

clean-build: install build
