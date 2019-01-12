package torrent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	slackTorrentService "github.com/eschizoid/flixctl/slack/torrent"
	"github.com/eschizoid/flixctl/torrent"
	"github.com/spf13/cobra"
)

var StatusTorrentCmd = &cobra.Command{
	Use:   "status",
	Short: "To Show Torrents Status",
	Long:  "to show the status of the torrents being downloaded",
	Run: func(cmd *cobra.Command, args []string) {
		Status()
	},
}

func Status() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	svc := ec2.New(awsSession, awsSession.Config)
	instanceID := ec2Service.FetchInstanceID(svc, awsResourceTagNameValue)
	if ec2status := ec2Service.Status(svc, instanceID); strings.EqualFold(ec2status, Ec2RunningStatus) {
		torrents := torrent.Status()
		body, _ := json.Marshal(torrents)
		fmt.Println(string(body))
		if notify, _ := strconv.ParseBool(slackNotification); notify {
			slackTorrentService.SendStatus(torrents, slackIncomingHookURL)
		}
	} else {
		m := make(map[string]string)
		m["plex_status"] = strings.ToLower(Ec2StoppedStatus)
		jsonString, _ := json.Marshal(m)
		fmt.Println("\n" + string(jsonString))
	}
}
