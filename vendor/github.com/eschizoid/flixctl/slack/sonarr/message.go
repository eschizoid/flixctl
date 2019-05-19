//nolint:dupl
package sonarr

import (
	"encoding/json"
	"fmt"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/jrudio/go-sonarr-client"
	"github.com/nlopes/slack"
)

func SendShows(results []sonarr.SearchResults, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, len(results))
	for _, sonarResult := range results {
		attachmentDateAdded := slack.AttachmentField{
			Title: "Date Added",
			Value: sonarResult.Added,
			Short: true,
		}
		attachment := slack.Attachment{
			Color:     "#5DBCD2",
			Title:     sonarResult.CleanTitle,
			TitleLink: sonarResult.Network,
			Fields: []slack.AttachmentField{
				attachmentDateAdded,
			},
			Footer:     "Sonarr Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/sonarr/7b3bc1604171f754.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		}
		attachments = append(attachments, attachment)
	}
	if len(attachments) == 0 {
		attachments = append(attachments, slack.Attachment{
			Color:      "#5DBCD2",
			Text:       "*No Shows found*",
			MarkdownIn: []string{"text"},
			Footer:     "Sonarr Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/sonarr/7b3bc1604171f754.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		})
	}
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending shows: [%s]\n", err)
	}
}

func SendShowsHelp(slackIncomingHookURL string) {
	const TitleLink = "https://github.com/eschizoid/flixctl/blob/master/README.adoc"

	attachmentRequestShows := slack.AttachmentField{
		Value: "✅ Request shows via Ombi:\n`/shows-request`",
		Short: false,
	}
	attachmentSearchShows := slack.AttachmentField{
		Value: "✅ Search shows using Sonarr client:\n`/shows-search`",
		Short: false,
	}
	attachment := slack.Attachment{
		Text: "👋 Need some help with `/shows-request` or `/shows-search`?",
		Fields: []slack.AttachmentField{
			attachmentRequestShows,
			attachmentSearchShows,
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
		fmt.Printf("Error while sending shows help: [%s]\n", err)
	}
}
