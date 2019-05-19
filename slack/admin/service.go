package admin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eschizoid/flixctl/models"
	"github.com/nlopes/slack"
)

func SaveToken(clientID string, token string, svc *dynamodb.DynamoDB) error {
	err := models.SaveOauthToken(clientID, token, svc)
	return err
}

func SendAdminHelp(slackIncomingHookURL string) {
	const TitleLink = "https://github.com/eschizoid/flixctl/blob/master/README.adoc"

	attachmentRenewCerts := slack.AttachmentField{
		Value: "✅ Renew TLS certificates for all plex related services:\n`/admin renew-certs`",
		Short: false,
	}
	attachmentRestartServices := slack.AttachmentField{
		Value: "✅ Restart all Plex related services:\n`/admin restart-services`",
		Short: false,
	}
	attachmentSlackToken := slack.AttachmentField{
		Value: "✅ Get a Slack oauth token for a given client id:\n`/admin slack-token`",
		Short: false,
	}
	attachmentSlackPurge := slack.AttachmentField{
		Value: "✅ Purge Slack messages from all channels:\n`/admin slack-purge`",
		Short: false,
	}
	attachment := slack.Attachment{
		Text: "👋 Need some help with `/admin`?",
		Fields: []slack.AttachmentField{
			attachmentRenewCerts,
			attachmentRestartServices,
			attachmentSlackToken,
			attachmentSlackPurge,
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
		fmt.Printf("Error while sending admin help: [%s]\n", err)
	}
}
