package plex

import (
	"fmt"

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
		fmt.Println("Plex Stopped")
	},
}

func Stop() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	svc := ec2.New(awsSession, awsSession.Config)
	var instanceID = ec2Service.FetchInstanceID(svc, "plex")
	if ec2Service.Status(svc, instanceID) == "Stopped" {
		slackService.SendStop(slackIncomingHookURL)
		return
	}
	var oldVolumeID = ebsService.FetchVolumeID(svc, "plex")
	snapService.Create(svc, oldVolumeID, "plex")
	ec2Service.Stop(svc, instanceID)
	ebsService.Detach(svc, oldVolumeID)
	ebsService.Delete(svc, oldVolumeID)
	if slackIncomingHookURL != "" {
		slackService.SendStop(slackIncomingHookURL)
	}
}
