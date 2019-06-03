package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/plex"
	slackService "github.com/eschizoid/flixctl/slack/plex"
)

func executePlexCommand(evt json.RawMessage) {
	var input models.Input
	fmt.Printf("Exectuing λ with payload: %+v\n", input)

	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Argument {
	case "enable-monitoring":
		plex.EnabledMonitorRule()
	case "disable-monitoring":
		plex.DisabledMonitorRule()
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
		slackService.SendPlexHelp(os.Getenv("SLACK_GENERAL_HOOK_URL"))
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	}
	fmt.Println("Successfully executed λ plex")
}

func main() {
	lambda.Start(executePlexCommand)
}
