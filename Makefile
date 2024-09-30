SHELL=/bin/bash
build:
	@go build -o tmp/bin/ cmd/main.go

run: build
	@./tmp/bin/main

watch:
	@air -c .air.toml
