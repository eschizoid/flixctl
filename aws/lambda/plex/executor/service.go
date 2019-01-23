package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/slack"
	"github.com/eschizoid/flixctl/cmd/plex"
)

func executePlexCommand(evt json.RawMessage) {
	var input slack.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Command {
	case "start":
		plex.Start()
	case "stop":
		plex.Stop()
	case "status":
		plex.Status()
	}
	fmt.Println("Successfully executed plex command")
}

func main() {
	lambda.Start(executePlexCommand)
}
