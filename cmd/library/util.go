package library

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/glacier"
	slackService "github.com/eschizoid/flixctl/slack/library"
)

func NotifySlack(jobDescriptions []glacier.JobDescription) {
	if notify, _ := strconv.ParseBool(slackNotification); notify {
		slackService.SendJobs(jobDescriptions, slackIncomingHookURL)
	}
}
