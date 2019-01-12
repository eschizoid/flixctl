package plex

import (
	"encoding/json"
	"fmt"
	"strconv"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ebsService "github.com/eschizoid/flixctl/aws/ebs"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	snapService "github.com/eschizoid/flixctl/aws/snapshot"
	slackService "github.com/eschizoid/flixctl/slack/plex"
	"github.com/spf13/cobra"
)

var StartPlexCmd = &cobra.Command{
	Use:   "start",
	Short: "To Start Plex",
	Long:  "to start the Plex Media Center.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		Start()
		close(shutdownCh)
		m := make(map[string]string)
		m["plex_status"] = "started"
		jsonString, _ := json.Marshal(m)
		fmt.Println("\n" + string(jsonString))
	},
}

func Start() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	svc := ec2.New(awsSession, awsSession.Config)
	var instanceID = ec2Service.FetchInstanceID(svc, awsResourceTagNameValue)
	var status = ec2Service.Status(svc, instanceID)
	if status == Ec2RunningStatus {
		slackService.SendStatus("running", slackIncomingHookURL)
		return
	}
	ec2Service.Start(svc, instanceID)
	var oldSnapshotID = snapService.FetchSnapshotID(svc, awsResourceTagNameValue)
	ebsService.Create(svc, oldSnapshotID, awsResourceTagNameValue)
	var newVolumeID = ebsService.FetchVolumeID(svc, awsResourceTagNameValue)
	ebsService.Attach(svc, instanceID, newVolumeID)
	snapService.Delete(svc, oldSnapshotID)
	if notify, _ := strconv.ParseBool(slackNotification); notify {
		slackService.SendStatus("running", slackIncomingHookURL)
	}
}
