.PHONY: build run 
build:
	@go build -o bin/lantern ./cmd/lantern
run: build
	@./bin/lantern
