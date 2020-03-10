GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get
BINARY_DIR=bin

all: test build
build:
	$(GOBUILD) -o ./$(BINARY_DIR)/ -v ./cmd/...
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f ./$(BINARY_DIR)/*
deps:
	$(GOMOD) tidy
