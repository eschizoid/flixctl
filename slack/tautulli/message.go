package tautulli

import (
	"encoding/json"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/nlopes/slack"
)

const (
	outgoingHookURL = "https://marianoflix.duckdns.org:9000/hooks/tautulli-notification"
)

func ForwardEvent(event string) error {
	var attachments []slack.Attachment
	attachments = append(attachments, slack.Attachment{
		Color:  "good",
		Text:   event,
		Footer: "Tautulli Server",
		Ts:     json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
	})
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(outgoingHookURL, message)
	return err
}
