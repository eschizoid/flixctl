package slack

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/eschizoid/flixctl/torrent"
	"github.com/nlopes/slack"
)

const (
	outgoingHookURL       = "https://marianoflix.duckdns.org:9000/hooks/download"
	incomingSearchHookURL = "https://hooks.slack.com/services/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"
)

func SendDownloadLinks(search *torrent.Search) {
	var attachments []slack.Attachment
	for _, torrentResult := range search.Out {
		encodedMagnetLink := base64.StdEncoding.EncodeToString([]byte(torrentResult.Magnet))
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
			Color:     "#36a64f",
			Title:     torrentResult.Name,
			TitleLink: outgoingHookURL + "?n=" + torrentResult.Name + "&m=" + encodedMagnetLink,
			Fields: []slack.AttachmentField{
				attachmentFieldSize,
				attachmentFieldSeeders,
				attachmentFieldUploadDate,
				attachmentFieldSource,
			},
		}
		attachments = append(attachments, attachment)
	}
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(incomingSearchHookURL, message)
	if err != nil {
		fmt.Printf("error while seding download links: %s\n", err)
	}
}

func NotifyDownloadStarted(torrentName string) {
	fmt.Print(torrentName)
}
