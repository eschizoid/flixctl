package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/slack"
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
		postToWebhooks(baseHookURL+slash.Command, map[string]interface{}{
			"filter": slash.Text,
		})
		fmt.Printf("Filter: %s", slash.Text)
		message = fmt.Sprintf(`{"response_type":"ephemeral", "text":"Executing library jobs command"}`)
	case "/library-initiate":
		postToWebhooks(baseHookURL+slash.Command, map[string]interface{}{})
		message = fmt.Sprintf(`{"response_type":"ephemeral", "text":"Executing library initiate command"}`)
	case "/library-catalogue":
		postToWebhooks(baseHookURL+slash.Command, map[string]interface{}{
			"filter": slash.Text,
		})
		message = fmt.Sprintf(`{"response_type":"ephemeral", "text":"Executing library catalogue command"}`)
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
	fullURL := url + fmt.Sprintf("?token=%s", slack.SigningSecret)
	fmt.Printf("Full url: %s", fullURL)
	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(byteMessage))
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
