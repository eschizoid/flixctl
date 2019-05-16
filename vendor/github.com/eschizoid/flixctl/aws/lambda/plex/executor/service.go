package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/plex"
	slackLambdaService "github.com/eschizoid/flixctl/slack/lambda"
)

func executePlexCommand(evt json.RawMessage) {
	var input models.Input
	fmt.Printf("Exectuing λ with payload: %+v\n", input)

	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Argument {
	case "enable-monitoring":
		activateMonitor := "true" //nolint:goconst
		plex.EnableDisableMonitor(activateMonitor)
	case "disable-monitoring":
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
	case "help":
		fmt.Printf("Executing %s command \n", input.Argument)
		slackLambdaService.SendPlexHelp("")
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	}
	fmt.Println("Successfully executed λ plex")
}

func main() {
	lambda.Start(executePlexCommand)
}
