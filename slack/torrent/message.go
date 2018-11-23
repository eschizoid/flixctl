package torrent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/eschizoid/flixctl/torrent"
	"github.com/hekmon/transmissionrpc"
	"github.com/nlopes/slack"
)

const (
	outgoingHookURL = "https://marianoflix.duckdns.org:9000/hooks/torrent-download"
)

func SendDownloadLinks(search *torrent.Search, slackIncomingHookURL string, directoryDir string) {
	var attachments []slack.Attachment
	token := os.Getenv("SLACK_MOVIES_SEARCH_TOKEN")
	if directoryDir == "" {
		directoryDir = os.TempDir()
	}
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
			Color: "#C40203",
			Title: torrentResult.Name,
			TitleLink: outgoingHookURL +
				"?n=" + url.QueryEscape(encodedName) +
				"&m=" + url.QueryEscape(encodedMagnetLink) +
				"&t=" + token +
				"&d=" + directoryDir,
			Fields: []slack.AttachmentField{
				attachmentFieldSize,
				attachmentFieldSeeders,
				attachmentFieldUploadDate,
				attachmentFieldSource,
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
	var attachments []slack.Attachment
	for _, torrentFile := range torrents {
		attachments = append(attachments, slack.Attachment{
			Color: "#C40203",
			Text: fmt.Sprintf(`
*Name*: %s
*Percentage*: %.2f%%
*ETA*: %d`, *torrentFile.Name, *torrentFile.PercentDone*100, *torrentFile.Eta/3600),
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
