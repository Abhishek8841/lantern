.PHONY: build run run-a run-b test vet clean

# override on the command line, e.g. make run LISTEN=:9000)
LISTEN   ?= :8080
IDENTITY ?= lantern-identity
CONNECT  ?=

FLAGS = --listen $(LISTEN) --identity $(IDENTITY)
ifneq ($(CONNECT),)
	FLAGS += --connect $(CONNECT)
endif

build:
	@go build -o bin/lantern ./cmd/lantern

run: build
	@./bin/lantern $(FLAGS)




# local two-peer testing
run-a: build
	@./bin/lantern --listen :8080 --identity lantern-a

run-b: build
	@./bin/lantern --listen :8081 --identity lantern-b

test:
	@go test ./...

vet:
	@go vet ./...

clean:
	@rm -rf bin/
