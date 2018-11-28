package plex

import (
	"strconv"

	slackService "github.com/eschizoid/flixctl/slack/plex"
)

func NotifySlack(status string) {
	if notify, _ := strconv.ParseBool(slackNotification); notify {
		slackService.SendStatus(status, slackIncomingHookURL)
	}
}
