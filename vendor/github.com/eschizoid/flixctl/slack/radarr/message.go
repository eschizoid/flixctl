package radarr

import (
	"encoding/json"
	"fmt"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/jrudio/go-radarr-client"
	"github.com/nlopes/slack"
)

func SendMovies(results []radarr.Movie, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, len(results))
	for _, radarResult := range results {
		attachmentDateAdded := slack.AttachmentField{
			Title: "Date Added",
			Value: radarResult.Added,
			Short: true,
		}
		attachment := slack.Attachment{
			Color:     "#C40203",
			Title:     radarResult.CleanTitle,
			TitleLink: radarResult.YouTubeTrailerID,
			Fields: []slack.AttachmentField{
				attachmentDateAdded,
			},
			Footer:     "Radarr Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/radarr/989ae3e6536a72dc.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		}
		attachments = append(attachments, attachment)
	}
	if len(attachments) == 0 {
		attachments = append(attachments, slack.Attachment{
			Color:      "#C40203",
			Text:       "*No Movies found*",
			MarkdownIn: []string{"text"},
			Footer:     "Torrent Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/radarr/989ae3e6536a72dc.png",
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
