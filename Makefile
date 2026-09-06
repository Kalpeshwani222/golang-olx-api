.PHONY: build run

build:
	@go build -o bin/api.exe ./cmd/api

run: build
	@./bin/api.exe