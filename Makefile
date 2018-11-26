# Go parameters
THIS_FILE := $(lastword $(MAKEFILE_LIST))
GOCMD := go
GOBUILD := $(GOCMD) build
GODEP :=dep
GOINSTALL := $(GOCMD) install
GOLINT := golangci-lint
GOFMT := gofmt
SHELL := /bin/bash
TARGET := $(shell echo $${PWD\#\#*/})
VERSION := 1.0.0
BUILD := `git rev-parse --short HEAD`
LDFLAGS=-ldflags "-X=main.VERSION=$(VERSION) -X=main.BUILD=$(BUILD)"
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL: $(TARGET)

all: lint install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: clean $(TARGET) build-lambdas
	@true

clean:
	@rm -f $(TARGET)
	@rm -rf $(shell pwd)/aws/lambda/plex/executor/executor
	@rm -rf $(shell pwd)/aws/lambda/plex/executor/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/plex/dispatcher/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/plex/dispatcher/dispatcher
	@rm -rf $(shell pwd)/aws/lambda/torrent/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/torrent/torrent
	@rm -rf $(shell pwd)/aws/lambda/library/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/library/library

install:
ifeq ($(UPDATE_VENDOR), true)
	@$(MAKE) -f $(THIS_FILE) update-vendor
endif
	$(GOINSTALL) $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	$(GOFMT) -l -w $(SRC)

simplify:
	$(GOFMT) -s -l -w $(SRC)

run: install
	@$(TARGET)

dep:
	$(GODEP) check
	$(GODEP) ensure -v

lint:
	$(GOLINT) -v --deadline=5m run --disable gochecknoglobals

update:
	$(GODEP) ensure -update -v

update-vendor:
	@cp -R aws/ vendor/github.com/eschizoid/flixctl/aws
	@cp -R cmd/ vendor/github.com/eschizoid/flixctl/cmd
	@cp -R library/ vendor/github.com/eschizoid/flixctl/library
	@cp -R models/ vendor/github.com/eschizoid/flixctl/models
	@cp -R slack/ vendor/github.com/eschizoid/flixctl/slack
	@cp -R torrent/ vendor/github.com/eschizoid/flixctl/torrent
	@cp -R worker/ vendor/github.com/eschizoid/flixctl/worker

build-lambdas: clean build-lambda-plex-dispatcher build-lambda-plex-executor build-lambda-torrent-router build-lambda-library-retriever

build-lambda-plex-dispatcher:
	@cd $(shell pwd)/aws/lambda/plex/dispatcher; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-plex-executor:
	@cd $(shell pwd)/aws/lambda/plex/executor; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-torrent-router:
	@cd $(shell pwd)/aws/lambda/torrent; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-library-retriever:
	@cd $(shell pwd)/aws/lambda/library; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

zip-lambdas: build-lambdas zip-lambda-plex-dispatcher zip-lambda-plex-executor zip-lambda-torrent-router zip-lambda-library-retriever

zip-lambda-plex-dispatcher:
	@cd $(shell pwd)/aws/lambda/plex/dispatcher; \
	zip -X lambda.zip dispatcher

zip-lambda-plex-executor:
	@cd $(shell pwd)/aws/lambda/plex/executor; \
	zip -X lambda.zip executor

zip-lambda-torrent-router:
	@cd $(shell pwd)/aws/lambda/torrent; \
	zip -X lambda.zip torrent

zip-lambda-library-retriever:
	@cd $(shell pwd)/aws/lambda/library; \
	zip -X lambda.zip library

deploy-lambdas: zip-lambdas deploy-lambda-plex-dispatcher deploy-lambda-plex-executor deploy-lambda-torrent-router

deploy-lambda-plex-dispatcher:
	@aws lambda update-function-code \
	--function-name plex \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/plex/dispatcher/lambda.zip

deploy-lambda-plex-executor:
	@aws lambda update-function-code \
	--function-name plex-command-executor \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/plex/executor/lambda.zip

deploy-lambda-torrent-router:
	@aws lambda update-function-code \
	--function-name torrent-router \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/torrent/lambda.zip

deploy-lambda-library-retriever:
	@aws lambda update-function-code \
	--function-name library-retiever \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/library/lambda.zip

tag:
	@git tag --force $(VERSION)
	@git push origin --tags --force
