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
			Color:     "#C40203",
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
			Color:      "#C40203",
			Text:       "*No Shows found*",
			MarkdownIn: []string{"text"},
			Footer:     "Torrent Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/sonarr/7b3bc1604171f754.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		})
	}
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending download links: [%s]\n", err)
	}
}
