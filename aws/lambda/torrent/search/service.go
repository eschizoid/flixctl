package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/torrent"
)

func executeTorrentSearchCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Command {
	case "movies-search", "shows-search":
		// search dynamo
		// search s3
		torrent.Search(input.Text)
	case "movies-status", "shows-status":
		torrent.Status()
	}
	fmt.Println("Successfully executed plex command")
}

func main() {
	lambda.Start(executeTorrentSearchCommand)
}
