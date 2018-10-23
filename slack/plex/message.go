package plex

import (
	"fmt"

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
	var attachments []slack.Attachment
	attachments = append(attachments, slack.Attachment{
		Color: "good",
		Text:  text,
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	return err
}
