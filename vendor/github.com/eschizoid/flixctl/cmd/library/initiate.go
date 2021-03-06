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

var InitiateLibraryCmd = &cobra.Command{
	Use:   "initiate",
	Short: "To Initiate Library Jobs",
	Long:  "to initiate library inventory retrieval or archive retrieval.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		initiateJobOutput := glacierService.InitiateJob(svc, archiveID)
		if notify, _ := strconv.ParseBool(slackNotification); notify {
			slackService.SendInitiatedJobNotification(initiateJobOutput, slackIncomingHookURL)
		}
		json, _ := json.Marshal(initiateJobOutput)
		fmt.Println(string(json))
	},
}
