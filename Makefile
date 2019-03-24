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
VERSION := 2.0.0
BUILD := `git rev-parse --short HEAD`
LDFLAGS=-ldflags "-X=main.VERSION=$(VERSION) -X=main.BUILD=$(BUILD)"
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

define environment
{
  "Variables": {
    "AWS_RESOURCE_TAG_NAME_VALUE": "$(AWS_RESOURCE_TAG_NAME_VALUE)",
    "ENABLE_LIBRARY_SYNC": "$(ENABLE_LIBRARY_SYNC)",
    "DYNAMODB_ENDPOINT": "$(DYNAMODB_ENDPOINT)",
    "FLIXCTL_HOST": "$(FLIXCTL_HOST)",
    "HOOKS_URL": "$(HOOKS_URL)",
    "PLEX_PASSWORD": "$(PLEX_PASSWORD)",
    "PLEX_TOKEN": "$(PLEX_TOKEN)",
    "PLEX_USER": "$(PLEX_USER)",
    "SLACK_CLIENT_ID": "$(SLACK_CLIENT_ID)",
    "SLACK_CLIENT_SECRET": "$(SLACK_CLIENT_SECRET)",
    "SLACK_REDIRECT_URI": "$(SLACK_REDIRECT_URI)",
    "SLACK_LEGACY_TOKEN": "$(SLACK_LEGACY_TOKEN)",
    "SLACK_MONITORING_HOOK_URL": "$(SLACK_MONITORING_HOOK_URL)",
    "SLACK_NEW_RELEASES_HOOK_URL": "$(SLACK_NEW_RELEASES_HOOK_URL)",
    "SLACK_NOTIFICATION": "$(SLACK_NOTIFICATION)",
    "SLACK_REQUESTS_HOOK_URL": "$(SLACK_REQUESTS_HOOK_URL)",
    "SLACK_SIGNING_SECRET": "$(SLACK_SIGNING_SECRET)",
    "TAUTULI_API_KEY": "$(TAUTULI_API_KEY)",
    "TR_AUTH": "$(TR_AUTH)",
    "UPDATE_VENDOR": "$(UPDATE_VENDOR)"
  }
}
endef
export environment

.DEFAULT_GOAL: $(TARGET)

all: lint install

$(TARGET): $(SRC)
	@go build $(LDFLAGS) -o $(TARGET)

build: clean $(TARGET) build-lambdas
	@true

create-deploy-directory:
	@mkdir deploy
	@cp flixctl deploy

clean:
	@rm -f $(TARGET)
	@rm -rf $(shell pwd)/aws/lambda/plex/dispatcher/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/plex/dispatcher/dispatcher
	@rm -rf $(shell pwd)/aws/lambda/plex/monitor/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/plex/monitor/monitor
	@rm -rf $(shell pwd)/aws/lambda/torrent/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/torrent/torrent
	@rm -rf $(shell pwd)/aws/lambda/library/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/library/library

env:
	@echo "$$environment"

install:
ifeq ($(UPDATE_VENDOR), true)
	@$(MAKE) -f $(THIS_FILE) update-vendor
endif
	@$(GOINSTALL) $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@$(GOFMT) -l -w $(SRC)

simplify:
	@$(GOFMT) -s -l -w $(SRC)

run: install
	@$(TARGET)

dep:
	@$(GODEP) check
	@$(GODEP) ensure -v

lint: fmt simplify
	@$(GOLINT) run -v \
	--deadline=5m \
	--disable gochecknoglobals \
	--disable lll \
	--max-same-issues 100

update:
	@$(GODEP) ensure -update -v

update-vendor:
	@cp -R aws/ vendor/github.com/eschizoid/flixctl/aws
	@cp -R cmd/ vendor/github.com/eschizoid/flixctl/cmd
	@cp -R library/ vendor/github.com/eschizoid/flixctl/library
	@cp -R models/ vendor/github.com/eschizoid/flixctl/models
	@cp -R slack/ vendor/github.com/eschizoid/flixctl/slack
	@cp -R torrent/ vendor/github.com/eschizoid/flixctl/torrent
	@cp -R worker/ vendor/github.com/eschizoid/flixctl/worker

build-lambdas: clean build-lambda-plex-dispatcher build-lambda-plex-monitor build-lambda-torrent-router build-lambda-library-router

build-lambda-plex-dispatcher:
	@cd $(shell pwd)/aws/lambda/plex/dispatcher; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-plex-monitor:
	@cd $(shell pwd)/aws/lambda/plex/monitor; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-torrent-router:
	@cd $(shell pwd)/aws/lambda/torrent; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-library-router:
	@cd $(shell pwd)/aws/lambda/library; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

zip-lambdas: build-lambdas zip-lambda-plex-dispatcher zip-lambda-plex-monitor zip-lambda-torrent-router zip-lambda-library-router

zip-lambda-plex-dispatcher:
	zip -j -X $(shell pwd)/aws/lambda/plex/dispatcher/lambda.zip \
	$(shell pwd)/aws/lambda/plex/dispatcher/dispatcher

zip-lambda-plex-monitor:
	zip -j -X $(shell pwd)/aws/lambda/plex/monitor/lambda.zip \
	$(shell pwd)/aws/lambda/plex/monitor/monitor

zip-lambda-torrent-router:
	zip -j -X $(shell pwd)/aws/lambda/torrent/lambda.zip \
	$(shell pwd)/aws/lambda/torrent/torrent

zip-lambda-library-router:
	zip -j -X $(shell pwd)/aws/lambda/library/lambda.zip \
	$(shell pwd)/aws/lambda/library/library

deploy-lambdas: zip-lambdas deploy-lambda-plex-dispatcher deploy-lambda-plex-monitor deploy-lambda-torrent-router deploy-lambda-library-router

deploy-lambda-plex-dispatcher:
	@aws lambda update-function-code \
	--function-name plex \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/plex/dispatcher/lambda.zip
	@aws lambda update-function-configuration \
	--function-name plex \
	--handler dispatcher \
	--region $(AWS_REGION) \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-plex-monitor:
	@aws lambda update-function-code \
	--function-name plex-monitor \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/plex/monitor/lambda.zip
	@aws lambda update-function-configuration \
	--function-name plex-monitor \
	--handler monitor \
	--region $(AWS_REGION) \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-torrent-router:
	@aws lambda update-function-code \
	--function-name torrent-router \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/torrent/lambda.zip
	@aws lambda update-function-configuration \
	--function-name torrent-router \
	--handler torrent \
	--region $(AWS_REGION) \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-library-router:
	@aws lambda update-function-code \
	--function-name library-router \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/library/lambda.zip
	@aws lambda update-function-configuration \
	--function-name library-router \
	--handler library \
	--region $(AWS_REGION) \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

tag:
	@git tag --force $(VERSION)
	@git push origin --tags --force
