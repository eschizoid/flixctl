//nolint:dupl
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
			Color:     "#FFC32E",
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
			Color:      "#FFC32E",
			Text:       "*No Movies found*",
			MarkdownIn: []string{"text"},
			Footer:     "Radarr Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/radarr/989ae3e6536a72dc.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		})
	}
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending movies: [%s]\n", err)
	}
}

func SendMoviesHelp(slackIncomingHookURL string) {
	const TitleLink = "https://github.com/eschizoid/flixctl/blob/master/README.adoc"

	attachmentRequestMovie := slack.AttachmentField{
		Value: "✅ Request movies via Ombi:\n`/movies-request`",
		Short: false,
	}
	attachmentSearchMovie := slack.AttachmentField{
		Value: "✅ Search movies using Radarr client:\n`/movies-search`",
		Short: false,
	}
	attachment := slack.Attachment{
		Text: "👋 Need some help with `/movies-request` or `/movies-search`?",
		Fields: []slack.AttachmentField{
			attachmentRequestMovie,
			attachmentSearchMovie,
		},
		MarkdownIn: []string{"text", "fields"},
		Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
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
		fmt.Printf("Error while sending movies help: [%s]\n", err)
	}
}
