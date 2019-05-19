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

func SendPlexHelp(slackIncomingHookURL string) {
	const TitleLink = "https://github.com/eschizoid/flixctl/blob/master/README.adoc"

	attachmentPlexDisableMonitoring := slack.AttachmentField{
		Value: "✅ Enable Plex monitoring:\n`/plex disable-monitoring`",
		Short: false,
	}
	attachmentPlexEnableMonitoring := slack.AttachmentField{
		Value: "✅ Disable Plex monitoring:\n`/plex enable-monitoring`",
		Short: false,
	}
	attachmentPlexStart := slack.AttachmentField{
		Value: "✅ Start Plex:\n`/plex start`",
		Short: false,
	}
	attachmentPlexStop := slack.AttachmentField{
		Value: "✅ Stop Plex:\n`/plex stop`",
		Short: false,
	}
	attachmentPlexStatus := slack.AttachmentField{
		Value: "✅ Get Plex status:\n`/plex status`",
		Short: false,
	}
	attachmentPlexToken := slack.AttachmentField{
		Value: "✅ Get Plex token:\n`/plex token`",
		Short: false,
	}
	attachment := slack.Attachment{
		Text: "👋 Need some help with `/plex`?",
		Fields: []slack.AttachmentField{
			attachmentPlexDisableMonitoring,
			attachmentPlexEnableMonitoring,
			attachmentPlexStart,
			attachmentPlexStop,
			attachmentPlexStatus,
			attachmentPlexToken,
		},
		MarkdownIn: []string{"text", "fields"},
	}
	attachmentLearnMore := slack.Attachment{
		Text: fmt.Sprintf("<http://%s|Learn More>", TitleLink),
	}
	message := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			attachment,
			attachmentLearnMore,
		},
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending plex help: [%s]\n", err)
	}
}
