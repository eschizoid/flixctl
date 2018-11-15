package torrent

import (
	"encoding/json"
	"fmt"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	slackPlexService "github.com/eschizoid/flixctl/slack/plex"
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
	var torrentStatus string
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	svc := ec2.New(awsSession, awsSession.Config)
	instanceID := ec2Service.FetchInstanceID(svc, "plex")
	ec2status := ec2Service.Status(svc, instanceID)
	torrents := torrent.Status(ec2status)
	if len(torrents) > 0 {
		body, err := json.Marshal(torrents)
		torrentStatus = string(body)
		if err != nil {
			fmt.Printf("Cannot parse response from transmission: [%s]\n", err)
			panic(err)
		}
		if slackIncomingHookURL != "" {
			slackTorrentService.SendStatus(torrents, slackIncomingHookURL)
		}
		fmt.Println(torrentStatus)
	} else {
		slackPlexService.SendStatus("stopped", slackIncomingHookURL)
		fmt.Println("Plex status stopped")
	}
}
