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
	@rm -rf $(shell pwd)/aws/lambda/library/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/library/library
	@rm -rf $(shell pwd)/aws/lambda/plex/executor/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/plex/executor/executor
	@rm -rf $(shell pwd)/aws/lambda/plex/monitor/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/plex/monitor/monitor
	@rm -rf $(shell pwd)/aws/lambda/slack/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/slack/slack
	@rm -rf $(shell pwd)/aws/lambda/torrent/batch/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/torrent/batch/batch
	@rm -rf $(shell pwd)/aws/lambda/torrent/download/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/torrent/download/download
	@rm -rf $(shell pwd)/aws/lambda/torrent/request/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/torrent/request/request
	@rm -rf $(shell pwd)/aws/lambda/torrent/search/lambda.zip
	@rm -rf $(shell pwd)/aws/lambda/torrent/search/search

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

build-lambdas: clean \
	build-lambda-library-executor \
	build-lambda-plex-executor \
	build-lambda-plex-monitor \
	build-lambda-slack-dispatcher \
	build-lambda-torrent-batch-executor \
	build-lambda-torrent-download-executor \
	build-lambda-torrent-request-executor \
	build-lambda-torrent-search-executor

build-lambda-library-executor:
	@cd $(shell pwd)/aws/lambda/library; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-plex-executor:
	@cd $(shell pwd)/aws/lambda/plex/executor; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-plex-monitor:
	@cd $(shell pwd)/aws/lambda/plex/monitor; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-slack-dispatcher:
	@cd $(shell pwd)/aws/lambda/slack; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-torrent-batch-executor:
	@cd $(shell pwd)/aws/lambda/torrent/batch; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-torrent-download-executor:
	@cd $(shell pwd)/aws/lambda/torrent/download; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-torrent-request-executor:
	@cd $(shell pwd)/aws/lambda/torrent/request; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

build-lambda-torrent-search-executor:
	@cd $(shell pwd)/aws/lambda/torrent/search; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)

zip-lambdas: build-lambdas \
	zip-lambda-library-executor \
	zip-lambda-plex-executor \
	zip-lambda-plex-monitor \
	zip-lambda-slack-dispatcher \
	zip-lambda-torrent-batch-executor \
	zip-lambda-torrent-download-executor \
	zip-lambda-torrent-request-executor \
	zip-lambda-torrent-search-executor

zip-lambda-library-executor:
	zip -j -X $(shell pwd)/aws/lambda/library/lambda.zip \
	$(shell pwd)/aws/lambda/library/library

zip-lambda-plex-executor:
	zip -j -X $(shell pwd)/aws/lambda/plex/executor/lambda.zip \
	$(shell pwd)/aws/lambda/plex/executor/executor

zip-lambda-plex-monitor:
	zip -j -X $(shell pwd)/aws/lambda/plex/monitor/lambda.zip \
	$(shell pwd)/aws/lambda/plex/monitor/monitor

zip-lambda-slack-dispatcher:
	zip -j -X $(shell pwd)/aws/lambda/slack/lambda.zip \
	$(shell pwd)/aws/lambda/slack/slack

zip-lambda-torrent-batch-executor:
	zip -j -X $(shell pwd)/aws/lambda/torrent/batch/lambda.zip \
	$(shell pwd)/aws/lambda/torrent/batch/batch

zip-lambda-torrent-download-executor:
	zip -j -X $(shell pwd)/aws/lambda/torrent/download/lambda.zip \
	$(shell pwd)/aws/lambda/torrent/download/download

zip-lambda-torrent-request-executor:
	zip -j -X $(shell pwd)/aws/lambda/torrent/request/lambda.zip \
	$(shell pwd)/aws/lambda/torrent/request/request

zip-lambda-torrent-search-executor:
	zip -j -X $(shell pwd)/aws/lambda/torrent/search/lambda.zip \
	$(shell pwd)/aws/lambda/torrent/search/search

deploy-lambdas: zip-lambdas \
	deploy-lambda-library-executor \
	deploy-lambda-plex-executor \
	deploy-lambda-plex-monitor \
	deploy-lambda-slack-dispatcher \
	deploy-lambda-torrent-batch-executor \
	deploy-lambda-torrent-download-executor \
	deploy-lambda-torrent-request-executor \
	deploy-lambda-torrent-search-executor

delete-lambdas:
	@aws lambda delete-function \
	--function-name library-executor
	@aws lambda delete-function \
	--function-name plex-executor
	@aws lambda delete-function \
	--function-name plex-monitor
	@aws lambda delete-function \
	--function-name slack-dispatcher
	@aws lambda delete-function \
	--function-name torrent-batch
	@aws lambda delete-function \
	--function-name torrent-download
	@aws lambda delete-function \
	--function-name torrent-request
	@aws lambda delete-function \
	--function-name torrent-search

deploy-lambda-library-executor:
	@aws lambda create-function \
	--function-name library-executor \
	--handler library \
	--region $(AWS_REGION) \
	--runtime go1.x \
	--role arn:aws:iam::623592657701:role/lambda_basic_execution \
	--zip-file fileb://$(shell pwd)/aws/lambda/library/lambda.zip \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-plex-executor:
	@aws lambda create-function \
	--function-name plex-executor \
	--handler executor \
	--region $(AWS_REGION) \
	--runtime go1.x \
	--role arn:aws:iam::623592657701:role/lambda_basic_execution \
	--zip-file fileb://$(shell pwd)/aws/lambda/plex/executor/lambda.zip \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-plex-monitor:
	@aws lambda create-function \
	--function-name plex-monitor \
	--handler monitor \
	--region $(AWS_REGION) \
	--runtime go1.x \
	--role arn:aws:iam::623592657701:role/lambda_basic_execution \
	--zip-file fileb://$(shell pwd)/aws/lambda/plex/monitor/lambda.zip \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-slack-dispatcher:
	@aws lambda create-function \
	--function-name slack-dispatcher \
	--handler slack \
	--region $(AWS_REGION) \
	--runtime go1.x \
	--role arn:aws:iam::623592657701:role/lambda_basic_execution \
	--zip-file fileb://$(shell pwd)/aws/lambda/slack/lambda.zip \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-torrent-batch-executor:
	@aws lambda create-function \
	--function-name torrent-batch \
	--handler batch \
	--region $(AWS_REGION) \
	--runtime go1.x \
	--role arn:aws:iam::623592657701:role/lambda_basic_execution \
	--zip-file fileb://$(shell pwd)/aws/lambda/torrent/batch/lambda.zip \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-torrent-download-executor:
	@aws lambda create-function \
	--function-name torrent-download \
	--handler download \
	--region $(AWS_REGION) \
	--runtime go1.x \
	--role arn:aws:iam::623592657701:role/lambda_basic_execution \
	--zip-file fileb://$(shell pwd)/aws/lambda/torrent/download/lambda.zip \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-torrent-request-executor:
	@aws lambda create-function \
	--function-name torrent-request \
	--handler request \
	--region $(AWS_REGION) \
	--runtime go1.x \
	--role arn:aws:iam::623592657701:role/lambda_basic_execution \
	--zip-file fileb://$(shell pwd)/aws/lambda/torrent/request/lambda.zip \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

deploy-lambda-torrent-search-executor:
	@aws lambda create-function \
	--function-name torrent-search \
	--handler search \
	--region $(AWS_REGION) \
	--runtime go1.x \
	--role arn:aws:iam::623592657701:role/lambda_basic_execution \
	--zip-file fileb://$(shell pwd)/aws/lambda/torrent/search/lambda.zip \
	--timeout 900 \
	--memory-size 128 \
	--runtime go1.x \
	--environment "$$environment"

tag:
	@git tag --force $(VERSION)
	@git push origin --tags --force
