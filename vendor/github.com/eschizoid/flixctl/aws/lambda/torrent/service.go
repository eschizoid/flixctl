package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/torrent"
	slackLambdaService "github.com/eschizoid/flixctl/slack/lambda"
)

func executeTorrentCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Command {
	case "torrent-search":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" { //nolint:goconst
			slackLambdaService.SendTorrentHelp("")
		} else {
			torrent.Search(input.Argument)
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "torrent-status":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" {
			slackLambdaService.SendTorrentHelp("")
		} else {
			torrent.Status()
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "torrent-movies-download":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" {
			slackLambdaService.SendTorrentHelp("")
		} else {
			torrent.Download(input.Argument, "/plex/movies")
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "torrent-shows-download":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" {
			slackLambdaService.SendTorrentHelp(os.Getenv("SLACK_GENERAL_HOOK_URL"))
		} else {
			torrent.Download(input.Argument, "/plex/shows")
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	}
	fmt.Println("Successfully executed λ torrent")
}

func main() {
	lambda.Start(executeTorrentCommand)
}
