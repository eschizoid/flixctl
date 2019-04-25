package main

import (
	"encoding/base64"
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
		torrent.Search(input.Argument)
	case "torrent-status":
		torrent.Status()
	case "torrent-movies-download":
		magnetLink, err := base64.StdEncoding.DecodeString(input.Argument)
		if err != nil {
			fmt.Println("Error decoding while decoding magnet link: ", err)
			panic(err)
		}
		torrent.Download(string(magnetLink), "/plex/movies")
	case "torrent-shows-download":
		magnetLink, err := base64.StdEncoding.DecodeString(input.Argument)
		if err != nil {
			fmt.Println("Error decoding while decoding magnet link: ", err)
			panic(err)
		}
		torrent.Download(string(magnetLink), "/plex/shows")
	}
	fmt.Println("Successfully executed plex command")
}

func main() {
	lambda.Start(executeTorrentCommand)
}
