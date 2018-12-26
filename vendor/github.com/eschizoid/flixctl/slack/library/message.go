package library

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/nlopes/slack"
)

func SendJobs(jobDescriptions []glacier.JobDescription, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, len(jobDescriptions))
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
					Type:  "button",
					Text:  "Start",
					URL:   fmt.Sprintf("https://marianoflix.duckdns.org:9091/hooks/retrieve-job?t=%s&i%s", *jobDescription.Action, *jobDescription.JobId),
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
	message := &slack.WebhookMessage{
		Text:        "*Library Jobs*",
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending library jobs: [%s]\n", err)
	}
}
