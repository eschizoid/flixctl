# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GODEP=dep
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
GOLINT=golangci-lint
GOTEST=$(GOCMD) test

BINARY:=flixctl
VERSION?=vlatest
PLATFORMS:=windows linux darwin
os=$(word 1, $@)

build: build-cli build-lambda-plex-dispatcher build-lambda-plex-executor build-lambda-torrent-router

build-cli:
	$(GOBUILD)

build-lambda-plex-dispatcher:
	cd $(shell pwd)/aws/lambda/plex/dispatcher; \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD)

build-lambda-plex-executor:
	cd $(shell pwd)/aws/lambda/plex/executor; \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD)

build-lambda-torrent-router:
	cd $(shell pwd)/aws/lambda/torrent; \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD)

clean:
	$(GOCLEAN); \
	rm -rf $(shell pwd)/aws/lambda/plex/executor/executor; \
	rm -rf $(shell pwd)/aws/lambda/plex/executor/lambda.zip; \
	rm -rf $(shell pwd)/aws/lambda/plex/dispatcher/lambda.zip; \
	rm -rf $(shell pwd)/aws/lambda/plex/dispatcher/dispatcher; \
	rm -rf $(shell pwd)/aws/lambda/torrent/lambda.zip; \
	rm -rf $(shell pwd)/aws/lambda/torrent/torrent;

deps:
	$(GODEP) check
	$(GODEP) ensure -v

lints:
	$(GOLINT) -v --skip-dirs='vendor' run

install:
	$(GOINSTALL)

update:
	$(GODEP) ensure -update

zip: clean build zip-lambda-plex-dispatcher zip-lambda-plex-executor zip-lambda-torrent-router

zip-lambda-plex-dispatcher:
	cd $(shell pwd)/aws/lambda/plex/dispatcher; \
	zip -X lambda.zip dispatcher

zip-lambda-plex-executor:
	cd $(shell pwd)/aws/lambda/plex/executor; \
	zip -X lambda.zip executor

zip-lambda-torrent-router:
	cd $(shell pwd)/aws/lambda/torrent; \
	zip -X lambda.zip torrent

deploy-lambdas:	clean build zip deploy-lambda-plex-dispatcher deploy-lambda-plex-executor deploy-lambda-torrent-router

deploy-lambda-plex-dispatcher:
	aws lambda update-function-code \
	--function-name plex \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/plex/dispatcher/lambda.zip

deploy-lambda-plex-executor:
	aws lambda update-function-code \
	--function-name plex-command-executor \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/plex/executor/lambda.zip

deploy-lambda-torrent-router:
	aws lambda update-function-code \
	--function-name torrent-router \
	--region $(AWS_REGION) \
	--zip-file fileb://$(shell pwd)/aws/lambda/torrent/lambda.zip

release-cli: $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64
