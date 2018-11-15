package plex

import (
	"strconv"

	slackService "github.com/eschizoid/flixctl/slack/plex"
)

func NotifySlack(status string) {
	notify, _ := strconv.ParseBool(slackNotification)
	if notify {
		slackService.SendStatus(status, slackIncomingHookURL)
	}
}
