package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/cmd/plex"
	"github.com/go-playground/form"
)

var commandRegexp = regexp.MustCompile(`^(start|stop|status)$`)

type Slash struct {
	Token       string `form:"token"`
	TeamID      string `form:"team_id"`
	TeamDomain  string `form:"team_domain"`
	ChannelID   string `form:"channel_id"`
	ChannelName string `form:"chann	el_name"`
	UserID      string `form:"user_id"`
	UserName    string `form:"user_name"`
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseURL string `form:"response_url"`
}

func router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return executePlexCommand(request)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func executePlexCommand(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	values, err := url.ParseQuery(request.Body)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}
	slash := new(Slash)
	err = form.NewDecoder().Decode(slash, values)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}
	if !commandRegexp.MatchString(slash.Text) || slash.Text == "" {
		return clientError(http.StatusBadRequest)
	}
	sendSlackResponse(slash.ResponseURL, slash.Text)
	//TODO figure out how to do this in an async way without making λ die
	switch slash.Text {
	case "start":
		plex.Start()
	case "stop":
		plex.Stop()
	case "status":
		plex.Status()
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       "",
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func sendSlackResponse(responseURL string, command string) {
	var jsonStr = []byte(fmt.Sprintf(`{"response_type": "in_channel", "text":"Executing %s command"}`, command))
	resp, err := http.Post(responseURL, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(body))
}

func main() {
	lambda.Start(router)
}
