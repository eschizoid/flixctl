package torrent

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/eschizoid/flixctl/torrent"
	"github.com/nlopes/slack"
)

const (
	outgoingHookURL = "https://marianoflix.duckdns.org:9000/hooks/torrent-download"
)

func SendDownloadLinks(search *torrent.Search, slackIncomingHookURL string) {
	var attachments []slack.Attachment
	for _, torrentResult := range search.Out {
		encodedMagnetLink := base64.StdEncoding.EncodeToString([]byte(torrentResult.Magnet))
		encodedName := base64.StdEncoding.EncodeToString([]byte(torrentResult.Name))
		attachmentFieldSize := slack.AttachmentField{
			Title: "Size",
			Value: torrentResult.Size,
			Short: true,
		}
		attachmentFieldSeeders := slack.AttachmentField{
			Title: "Number of Seeders",
			Value: strconv.Itoa(torrentResult.Seeders),
			Short: true,
		}
		attachmentFieldUploadDate := slack.AttachmentField{
			Title: "Upload Date",
			Value: torrentResult.UploadDate,
			Short: true,
		}
		attachmentFieldSource := slack.AttachmentField{
			Title: "Source",
			Value: torrentResult.Source,
			Short: true,
		}
		attachment := slack.Attachment{
			Color: "good",
			Title: torrentResult.Name,
			TitleLink: outgoingHookURL +
				"?n=" + url.QueryEscape(encodedName) +
				"&m=" + url.QueryEscape(encodedMagnetLink) +
				"&t=" + os.Getenv("SLACK_SEARCH_TOKEN"),
			Fields: []slack.AttachmentField{
				attachmentFieldSize,
				attachmentFieldSeeders,
				attachmentFieldUploadDate,
				attachmentFieldSource,
			},
		}
		attachments = append(attachments, attachment)
	}
	if len(attachments) == 0 {
		attachments = append(attachments, slack.Attachment{
			Color:      "warning",
			Text:       "*No Torrents found*",
			MarkdownIn: []string{"text"},
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

func SendDownloadStart(envTorrentName string, slackIncomingHookURL string) {
	decodedTorrentName, err := base64.StdEncoding.DecodeString(envTorrentName)
	if err != nil {
		fmt.Printf("Could not decode torrent name: [%s]\n", err)
	}
	var attachments []slack.Attachment
	attachments = append(attachments, slack.Attachment{
		Color: "good",
		Text:  fmt.Sprintf("Starting to download %s!", string(decodedTorrentName)),
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err = slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending download start notification: [%s]\n", err)
	}
}

func SendStatus(status string, slackIncomingHookURL string) {
	var attachments []slack.Attachment
	var color string
	if status == "Command timed out" || status == "Plex Stopped" {
		color = "warning"
	} else {
		color = "good"
	}
	attachments = append(attachments, slack.Attachment{
		Color:      color,
		Text:       "```" + fmt.Sprint(status) + "```",
		MarkdownIn: []string{"text"},
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending torrents download status: [%s]\n", err)
	}
}
