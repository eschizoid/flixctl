package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/slack"
	"github.com/eschizoid/flixctl/cmd/plex"
	"github.com/go-playground/form"
)

var (
	baseHookURL = fmt.Sprintf("https://%s:9000/hooks", os.Getenv("FLIXCTL_HOST"))
)

func router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return dispatch(request)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func dispatch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { //nolint:gocyclo
	var message string
	if plexStatus := plex.Status(); strings.EqualFold(plexStatus, "Running") {
		values, err := url.ParseQuery(request.Body)
		if err != nil {
			return clientError(http.StatusBadRequest)
		}
		slash := new(slack.Slash)
		err = form.NewDecoder().Decode(slash, values)
		if err != nil {
			return clientError(http.StatusUnprocessableEntity)
		}
		if !slack.VerifySlackRequest(request) {
			return clientError(http.StatusForbidden)
		}
		switch slashCommand := slash.Command; slashCommand {
		case "/library-jobs":
			token := os.Getenv("SLACK_STATUS_TOKEN")
			if slash.Token != token {
				return clientError(http.StatusForbidden)
			}
			postToWebhooks(baseHookURL+slash.Command, map[string]interface{}{
				"token":  token,
				"filter": slash.Text,
			})
			message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Executing movies jobs command"}`)
		case "/library-initiate":
			message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Executing movies initiate command"}`)
		case "/torrent-catalogue":
			message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Executing catalogue command"}`)
		case "/torrent-download":
			message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Executing download command"}`)
		case "/torrent-upload":
			message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Executing upload command"}`)
		}
	} else {
		message = fmt.Sprintf(`{"response_type":"in_channel", "text":"Make sure Plex is running"}`)
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
		fmt.Printf("Error while sending post to webhooks: [%s]\n", err)
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
