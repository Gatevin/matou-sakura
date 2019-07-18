# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

CURRENT_DIR=$(shell pwd)
LGOPATH=$(shell echo ${CURRENT_DIR}/../..)

BINARY_NAME=matou-sakura-collector

release: clean release_build

release_build:
	export GOPATH=$(LGOPATH);$(GOBUILD) -o $(CURRENT_DIR)/bin/$(BINARY_NAME) -v main.go

clean:
	$(GOCLEAN)
