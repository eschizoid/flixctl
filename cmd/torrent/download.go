package torrent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	slackService "github.com/eschizoid/flixctl/slack/torrent"
	torrentService "github.com/eschizoid/flixctl/torrent"
	"github.com/spf13/cobra"
)

var DownloadTorrentCmd = &cobra.Command{
	Use:   "download",
	Short: "To Download A Torrent",
	Long:  "to download a torrent using Transmission client.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := ec2.New(awsSession, awsSession.Config)
		instanceID := ec2Service.FetchInstanceID(svc, "plex")
		if ec2status := ec2Service.Status(svc, instanceID); strings.EqualFold(ec2status, Ec2RunningStatus) {
			torrent := torrentService.TriggerDownload(magnetLink, downloadDir)
			body, _ := json.Marshal(torrent)
			fmt.Println(string(body))
			if notify, _ := strconv.ParseBool(slackNotification); notify {
				slackService.SendDownloadStart(*torrent.Name, slackIncomingHookURL)
			}
		} else {
			m := make(map[string]string)
			m["plex_status"] = strings.ToLower(Ec2StoppedStatus)
			jsonString, _ := json.Marshal(m)
			fmt.Println("\n" + string(jsonString))
		}
	},
}
