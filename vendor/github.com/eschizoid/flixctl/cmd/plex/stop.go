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

var StopPlexCmd = &cobra.Command{
	Use:   "stop",
	Short: "To Stop Plex",
	Long:  "to stop the Plex Media Center.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		Stop()
		close(shutdownCh)
		m := make(map[string]string)
		m["plex_status"] = "stopped"
		jsonString, _ := json.Marshal(m)
		fmt.Println("\n" + string(jsonString))
	},
}

func Stop() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	svc := ec2.New(awsSession, awsSession.Config)
	var instanceID = ec2Service.FetchInstanceID(svc, awsResourceTagNameValue)
	var status = ec2Service.Status(svc, instanceID)
	if status == Ec2StoppedStatus {
		slackService.SendStatus("stopped", slackIncomingHookURL)
		return
	}
	var oldVolumeID = ebsService.FetchVolumeID(svc, awsResourceTagNameValue)
	snapService.Create(svc, oldVolumeID, awsResourceTagNameValue)
	ec2Service.Stop(svc, instanceID)
	ebsService.Detach(svc, oldVolumeID)
	ebsService.Delete(svc, oldVolumeID)
	if notify, _ := strconv.ParseBool(slackNotification); notify {
		slackService.SendStatus("stopped", slackIncomingHookURL)
	}
}
