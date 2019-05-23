package torrent

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/eschizoid/flixctl/torrent"
	"github.com/hekmon/transmissionrpc"
	"github.com/nlopes/slack"
)

func SendDownloadLinks(search *torrent.Search, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, len(search.Out))
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
			Value: url.QueryEscape(torrentResult.Magnet),
			Short: false,
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
	var attachments = make([]slack.Attachment, 0, len(torrents))
	for _, torrentFile := range torrents {
		var name = fmt.Sprintf("*Name*: %s", *torrentFile.Name)
		var percentDone = fmt.Sprintf("*Percentage completed*: %.2f%%", *torrentFile.PercentDone*100)
		var eta = fmt.Sprintf("*ETA*: ~%.2f minutes", float64(*torrentFile.Eta)/60)
		var totalSize = fmt.Sprintf("*Size*: %.2f GB", float64(*torrentFile.TotalSize)/8e9)
		attachments = append(attachments, slack.Attachment{
			Color: "#C40203",
			Text: fmt.Sprintf(`%s
%s
%s
%s`, name, percentDone, totalSize, eta),
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

func SendTorrentHelp(slackIncomingHookURL string) {
	const TitleLink = "https://github.com/eschizoid/flixctl/blob/master/README.adoc"

	attachmentTorrentSearch := slack.AttachmentField{
		Value: "✅ Search for a movie using the given keyword(s):\n`/torrent-search`",
		Short: false,
	}
	attachmentTorrentStatus := slack.AttachmentField{
		Value: "✅ Show the status the shows and movies being downloaded:\n`/torrent-status`",
		Short: false,
	}
	attachmentTorrentMoviesDownload := slack.AttachmentField{
		Value: "✅ Download a movie using Transmission client:\n`/torrent-movies-download`",
		Short: false,
	}
	attachmentTorrentShowsDownload := slack.AttachmentField{
		Value: "✅ Download a show using Transmission client:\n`/torrent-shows-download`",
		Short: false,
	}
	attachment := slack.Attachment{
		Text: "👋 Need some help with `/torrent-search`, `/torrent-status`, `/torrent-movies-download` or `/torrent-shows-download`?",
		Fields: []slack.AttachmentField{
			attachmentTorrentSearch,
			attachmentTorrentStatus,
			attachmentTorrentMoviesDownload,
			attachmentTorrentShowsDownload,
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
		fmt.Printf("Error while sending torrent help: [%s]\n", err)
	}
}
