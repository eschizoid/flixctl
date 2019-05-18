package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/admin"
	slackLambdaService "github.com/eschizoid/flixctl/slack/lambda"
)

func executeAdminCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	fmt.Printf("Exectuing λ with payload: %+v\n", input)
	switch input.Argument {
	case "renew-certs":
		fmt.Printf("Executing %s command \n", input.Argument)
		admin.RenewCerts()
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "restart-services":
		fmt.Printf("Executing %s command \n", input.Argument)
		admin.RestartPlexServices()
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "slack-token":
		fmt.Printf("Executing %s command \n", input.Argument)
		admin.SlackOauthToken()
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "slack-purge":
		fmt.Printf("Executing %s command \n", input.Argument)
		admin.SlackPurge()
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "help":
		fmt.Printf("Executing %s command \n", input.Argument)
		slackLambdaService.SendAdminHelp(os.Getenv("SLACK_GENERAL_HOOK_URL"))
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	}
}

func main() {
	lambda.Start(executeAdminCommand)
}
