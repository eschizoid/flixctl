package library

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/eschizoid/flixctl/models"
	util "github.com/eschizoid/flixctl/slack"
	"github.com/nlopes/slack"
)

func SendJobs(jobDescriptions []*glacier.JobDescription, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, len(jobDescriptions))
	token := util.SigningSecret
	for _, jobDescription := range jobDescriptions {
		var url string
		if action := *jobDescription.Action; action == "InventoryRetrieval" {
			url = util.LibraryInventoryHookURL + "?i=" + *jobDescription.JobId + "&e=" + "true" + "&token=" + token
		} else {
			url = util.LibraryDownloadHookURL + "?i=" + *jobDescription.JobId + "&token=" + token
		}
		attachmentFieldJobType := slack.AttachmentField{
			Title: "*Job Type*",
			Value: *jobDescription.Action,
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
				attachmentFieldCreatedAt,
				attachmentFieldStatusCode,
				attachmentFieldCompleted,
			},
			MarkdownIn: []string{"text", "fields"},
			Actions: []slack.AttachmentAction{
				{
					Type:  "button",
					Text:  "Download",
					URL:   url,
					Style: "default",
					Confirm: &slack.ConfirmationField{
						Title:       "Are you sure you want to start the job?",
						Text:        "Is the job completed?",
						OkText:      "Yes",
						DismissText: "No",
					},
				},
			},
		})
	}
	if len(jobDescriptions) == 0 {
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

func SendInventory(archives []models.InventoryArchive, slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, len(archives))
	token := util.SigningSecret
	for _, archive := range archives {
		attachmentTitle := slack.AttachmentField{
			Title: "*Archive Description*",
			Value: archive.ArchiveDescription,
			Short: true,
		}
		attachmentSize := slack.AttachmentField{
			Title: "*Size*",
			Value: strconv.Itoa(archive.Size),
			Short: true,
		}
		attachmentUploadDate := slack.AttachmentField{
			Title: "*Upload Date*",
			Value: archive.CreationDate,
			Short: true,
		}
		attachments = append(attachments, slack.Attachment{
			Color: "#C97D27",
			Fields: []slack.AttachmentField{
				attachmentTitle,
				attachmentSize,
				attachmentUploadDate,
			},
			MarkdownIn: []string{"text", "fields"},
			Actions: []slack.AttachmentAction{
				{
					Type:  "button",
					Text:  "Start",
					Style: "default",
					URL: util.LibraryInitiateArchiveHookURL +
						"?i=" + archive.ArchiveID +
						"&token=" + token,
					Confirm: &slack.ConfirmationField{
						Title:       "Are you sure you want to start the job?",
						OkText:      "Yes",
						DismissText: "No",
					},
				},
			},
		})
	}
	if len(archives) == 0 {
		attachments = append(attachments, slack.Attachment{
			Color:      "#C97D27",
			Text:       "*No Library Inventory Found*",
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
		fmt.Printf("Error while sending library inventory: [%s]\n", err)
	}
}

func SendCatalogue(archives []models.Movie, slackIncomingHookURL string) {

}

func SendInitiatedJobNotification(job *glacier.InitiateJobOutput, slackIncomingHookURL string) {

}
