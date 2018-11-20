package library

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/glacier"
	slackService "github.com/eschizoid/flixctl/slack/library"
)

func NotifySlack(jobDescriptions []*glacier.JobDescription) {
	notify, _ := strconv.ParseBool(slackNotification)
	if notify {
		slackService.SendJobs(jobDescriptions, slackIncomingHookURL)
	}
}
