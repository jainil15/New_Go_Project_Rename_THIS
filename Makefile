SHELL=/bin/bash
build:
	@go build -o tmp/build/ cmd/main.go

run: build
	./tmp/build/main
