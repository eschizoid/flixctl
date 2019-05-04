package plex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/nlopes/slack"
)

func SendStatus(status string, slackIncomingHookURL string) {
	var attachments []slack.Attachment
	attachments = append(attachments, slack.Attachment{
		Color:      "#C97D27",
		Text:       fmt.Sprintf("Plex status: *%s*", strings.ToLower(status)),
		Footer:     "Plex Server",
		FooterIcon: "https://emoji.slack-edge.com/TD00VE755/plex/a1379540fa1021c2.png",
		MarkdownIn: []string{"text"},
		Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending plex status: [%s]\n", err)
	}
}

func SendToken(token string, slackIncomingHookURL string) {
	var attachments []slack.Attachment
	attachments = append(attachments, slack.Attachment{
		Color:      "#C97D27",
		Text:       fmt.Sprintf("Plex token: *%s*", token),
		Footer:     "Plex Server",
		FooterIcon: "https://emoji.slack-edge.com/TD00VE755/plex/a1379540fa1021c2.png",
		MarkdownIn: []string{"text"},
		Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending plex token: [%s]\n", err)
	}
}
