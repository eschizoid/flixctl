package plex

import (
	"encoding/json"
	"fmt"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/nlopes/slack"
)

func SendStart(slackIncomingHookURL string) {
	err := sendMessage(fmt.Sprint("Plex successfully started!"), slackIncomingHookURL)
	if err != nil {
		fmt.Printf("Error while sending plex start notification: [%s]\n", err)
	}
}

func SendStop(slackIncomingHookURL string) {
	err := sendMessage(fmt.Sprint("Plex successfully stopped!"), slackIncomingHookURL)
	if err != nil {
		fmt.Printf("Error while sending plex stop notification: [%s]\n", err)
	}
}

func SendStatus(status string, slackIncomingHookURL string) {
	err := sendMessage(fmt.Sprintf("Plex status: %s", status), slackIncomingHookURL)
	if err != nil {
		fmt.Printf("Error while sending plex status: [%s]\n", err)
	}
}

func sendMessage(text string, slackIncomingHookURL string) error {
	println(text)
	var attachments []slack.Attachment
	attachments = append(attachments, slack.Attachment{
		Color:  "good",
		Text:   text,
		Footer: "Plex Server",
		Ts:     json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	return err
}
