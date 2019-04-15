package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/torrent"
)

func executeTorrentCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Command {
	case "torrent-search":
		torrent.Search(input.Parameter)
	case "torrent-status":
		torrent.Status()
	case "torrent-movies-download":
		torrent.Download(input.Parameter, "/plex/movies")
	case "torrent-shows-download":
		torrent.Download(input.Parameter, "/plex/shows")
	}
	fmt.Println("Successfully executed plex command")
}

func main() {
	lambda.Start(executeTorrentCommand)
}
