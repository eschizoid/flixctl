package torrent

import (
	"encoding/json"
	"fmt"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/eschizoid/flixctl/torrent"
	"github.com/hekmon/transmissionrpc"
	"github.com/nlopes/slack"
)

func SendDownloadLinks(search *torrent.Search, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, len(search.Out))
	for _, torrentResult := range search.Out {
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
		attachmentFieldMagnetLink := slack.AttachmentField{
			Title: "Magnet Link",
			Value: torrentResult.Magnet,
			Short: true,
		}
		attachment := slack.Attachment{
			Color:     "#C40203",
			Title:     torrentResult.Name,
			TitleLink: util.TorrentDownloadHookURL,
			Fields: []slack.AttachmentField{
				attachmentFieldSize,
				attachmentFieldSeeders,
				attachmentFieldUploadDate,
				attachmentFieldSource,
				attachmentFieldMagnetLink,
			},
			Footer:     "Torrent Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/transmission/51fa8bddc5425861.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		}
		attachments = append(attachments, attachment)
	}
	if len(attachments) == 0 {
		attachments = append(attachments, slack.Attachment{
			Color:      "#C40203",
			Text:       "*No Torrents found*",
			MarkdownIn: []string{"text"},
			Footer:     "Torrent Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/transmission/51fa8bddc5425861.png",
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

func SendDownloadStart(torrentName string, slackIncomingHookURL string) {
	var attachments []slack.Attachment
	attachments = append(attachments, slack.Attachment{
		Color:      "#C40203",
		Text:       fmt.Sprintf("Starting to download *%s*!", torrentName),
		MarkdownIn: []string{"text"},
		Footer:     "Torrent Client",
		FooterIcon: "https://emoji.slack-edge.com/TD00VE755/transmission/51fa8bddc5425861.png",
		Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending download start notification: [%s]\n", err)
	}
}

func SendStatus(torrents []transmissionrpc.Torrent, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, len(torrents))
	for _, torrentFile := range torrents {
		var name = fmt.Sprintf("*Name*: %s", *torrentFile.Name)
		var percentDone = fmt.Sprintf("*Percentage*: %.2f%%", *torrentFile.PercentDone*100)
		var eta = fmt.Sprintf("*ETA*: %.2f minutes", float64(*torrentFile.Eta/60))
		attachments = append(attachments, slack.Attachment{
			Color: "#C40203",
			Text: fmt.Sprintf(`%s
%s
%s`, name, percentDone, eta),
			MarkdownIn: []string{"text"},
			Footer:     "Torrent Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/transmission/51fa8bddc5425861.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		})
	}
	if len(attachments) == 0 {
		attachments = append(attachments, slack.Attachment{
			Color:      "#C40203",
			Text:       "*No Torrents found*",
			MarkdownIn: []string{"text"},
			Footer:     "Torrent Client",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/transmission/51fa8bddc5425861.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		})
	}
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending torrents download status: [%s]\n", err)
	}
}
