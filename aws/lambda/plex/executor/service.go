package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/plex"
)

func executePlexCommand(evt json.RawMessage) {
	var input models.Input
	fmt.Printf("Exectuing λ with payload: %+v\n", input)

	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Argument {
	case "enable-monitor":
		activateMonitor := "true" //nolint:goconst
		plex.EnableDisableMonitor(activateMonitor)
	case "disable-monitor":
		activateMonitor := "false" //nolint:goconst
		plex.EnableDisableMonitor(activateMonitor)
	case "start":
		plex.Start()
	case "stop":
		slackNotification := "true" //nolint:goconst
		plex.Stop(slackNotification)
	case "status":
		plex.Status()
	case "token":
		plex.Token()
	}
	fmt.Println("Successfully executed plex command")
}

func main() {
	lambda.Start(executePlexCommand)
}
