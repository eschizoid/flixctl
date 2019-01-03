package library

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/service/glacier"
	util "github.com/eschizoid/flixctl/slack"
	"github.com/nlopes/slack"
)

func SendJobs(jobDescriptions []*glacier.JobDescription, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, len(jobDescriptions))
	token := os.Getenv("SLACK_MOVIES_SEARCH_TOKEN")
	for _, jobDescription := range jobDescriptions {
		attachmentFieldJobType := slack.AttachmentField{
			Title: "*Job Type*",
			Value: *jobDescription.Action,
			Short: true,
		}
		attachmentFieldJobDescription := slack.AttachmentField{
			Title: "*Job Description*",
			Value: *jobDescription.JobDescription,
			Short: true,
		}
		attachmentFieldCreatedAt := slack.AttachmentField{
			Title: "*Creation Date*",
			Value: *jobDescription.CreationDate,
			Short: true,
		}
		attachmentFieldStatusCode := slack.AttachmentField{
			Title: "Status",
			Value: *jobDescription.StatusCode,
			Short: true,
		}
		attachmentFieldCompleted := slack.AttachmentField{
			Title: "*Completed*",
			Value: strconv.FormatBool(*jobDescription.Completed),
			Short: true,
		}
		attachments = append(attachments, slack.Attachment{
			Color: "#C97D27",
			Fields: []slack.AttachmentField{
				attachmentFieldJobType,
				attachmentFieldJobDescription,
				attachmentFieldCreatedAt,
				attachmentFieldStatusCode,
				attachmentFieldCompleted,
			},
			MarkdownIn: []string{"text", "fields"},
			Actions: []slack.AttachmentAction{
				{
					Type: "button",
					Text: "Start",
					URL: util.LibraryInventoryHookURL +
						"?i=" + *jobDescription.JobId +
						"&s=" + "true" +
						"&token=" + token,
					Style: "default",
					Confirm: &slack.ConfirmationField{
						Title:       "Are you sure you want to start the job retrieval?",
						Text:        "Is the job completed?",
						OkText:      "Yes",
						DismissText: "No",
					},
				},
			},
		})
	}
	if len(attachments) == 0 {
		attachments = append(attachments, slack.Attachment{
			Color:      "#C97D27",
			Text:       "*No Library Jobs Found*",
			MarkdownIn: []string{"text"},
			Footer:     "Plex Server",
			FooterIcon: "https://emoji.slack-edge.com/TD00VE755/plex/a1379540fa1021c2.png",
			Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
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
