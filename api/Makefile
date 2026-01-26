# Makefile

BINARY_NAME=api

# Default url: http://localhost:7000

build:
	go mod tidy && \
   	templ generate && \
	go generate && \
	go build -ldflags="-w -s" -o ./bin/main.exe

lint:
	golangci-lint run 

dev: 
	wgo run ./cmd/api

reset:
	go run ./cmd/resetter

# 	make -j2 lint