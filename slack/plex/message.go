package plex

import (
	"fmt"

	"github.com/nlopes/slack"
)

const (
	plexIncomingHookURL = "https://hooks.slack.com/services/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"
)

func SendStart() {
	err := sendMessage(fmt.Sprint("Plex successfully started!"))
	if err != nil {
		fmt.Printf("Error while sending plex start notification: [%s]\n", err)
	}
}

func SendStop() {
	err := sendMessage(fmt.Sprint("Plex successfully stopped!"))
	if err != nil {
		fmt.Printf("Error while sending plex stop notification: [%s]\n", err)
	}
}

func SendStatus(status string) {
	err := sendMessage(fmt.Sprintf("Plex status: %s", status))
	if err != nil {
		fmt.Printf("Error while sending plex status: [%s]\n", err)
	}
}

func sendMessage(text string) error {
	var attachments []slack.Attachment
	attachments = append(attachments, slack.Attachment{
		Color: "good",
		Text:  text,
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(plexIncomingHookURL, message)
	return err
}
