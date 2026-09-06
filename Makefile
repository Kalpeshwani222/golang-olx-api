.PHONY: build run

build:
	@go build -o bin/api.exe ./cmd/api

run: build
	@./bin/api.exe

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down
