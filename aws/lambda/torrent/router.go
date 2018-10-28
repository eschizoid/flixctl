package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/cmd/plex"
)

const (
	baseHookURL = "https://marianoflix.duckdns.org:9000/hooks"
)

func router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var message string
	switch request.HTTPMethod {
	case "POST":
		plexStatus := plex.Status()
		// Send request to webhooks
		if plexStatus == "Running" {
			path := request.Path
			if path == "/torrent-search" {
				postToWebhooks(baseHookURL+path, map[string]interface{}{"text": request.Body})
				message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Executing search command"}`)
			} else if path == "/torrent-status" {
				postToWebhooks(baseHookURL+path, map[string]interface{}{})
				message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Executing status command"}`)
			}
		} else {
			message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Please start Plex"}`)
		}
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       message,
	}, nil
}

func postToWebhooks(url string, message map[string]interface{}) {
	byteMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Unable to parse the request: [%s]\n", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(byteMessage))
	if err != nil {
		fmt.Printf("Error While sending post to webhooks: [%s]\n", err)
	}
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Unable to parse the response [%s]\n", err)
	}
	fmt.Println(result)
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
