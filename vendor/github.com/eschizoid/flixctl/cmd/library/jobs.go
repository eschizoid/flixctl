package library

import (
	"encoding/json"
	"fmt"
	"strconv"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	slackService "github.com/eschizoid/flixctl/slack/library"
	"github.com/spf13/cobra"
)

var JobsLibraryCmd = &cobra.Command{
	Use:   "jobs",
	Short: "To List Library Jobs",
	Long:  "to list library jobs.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		jobList := glacierService.ListJobs(svc)
		filteredJobs := chooseSucceededJobs(jobList.JobList, func(retrieval string) bool {
			return retrieval == retrievalType
		})
		json, _ := json.Marshal(filteredJobs)
		if notify, _ := strconv.ParseBool(slackNotification); notify {
			slackService.SendJobs(filteredJobs, slackIncomingHookURL)
		}
		fmt.Println(string(json))
	},
}

func chooseSucceededJobs(jobDescriptions []*glacier.JobDescription, test func(string) bool) (filteredJobs []glacier.JobDescription) {
	for _, job := range jobDescriptions {
		if test(*job.Action) {
			filteredJobs = append(filteredJobs, *job)
		}
	}
	return
}
