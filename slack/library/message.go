package library

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/glacier"
	util "github.com/eschizoid/flixctl/slack"
	"github.com/nlopes/slack"
)

func SendJobs(jobDescriptions []*glacier.JobDescription, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, len(jobDescriptions))
	for _, jobDescription := range jobDescriptions {
		var url string
		if strings.EqualFold(*jobDescription.Action, "InventoryRetrieval") {
			url = fmt.Sprintf("https://marianoflix.duckdns.org:9091/hooks/%s?t=%s&j%s", "retrieve-inventory", *jobDescription.Action, *jobDescription.JobId)
		} else {
			url = fmt.Sprintf("https://marianoflix.duckdns.org:9091/hooks/%s?t=%s&j%s", "retrieve-archive", *jobDescription.Action, *jobDescription.JobId)
		}
		attachments = append(attachments, slack.Attachment{
			Color:      "#C97D27",
			Text:       "Current Jobs",
			Footer:     "Plex Server",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/plex/a1379540fa1021c2.png",
			MarkdownIn: []string{"text"},
			Actions: []slack.AttachmentAction{
				{
					Type:  "button",
					Text:  "Start",
					URL:   url,
					Style: "primary",
				},
			},
			Ts: json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
		})
	}
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending library jobs: [%s]\n", err)
	}
}
