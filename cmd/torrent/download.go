package torrent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	plexTorrentService "github.com/eschizoid/flixctl/slack/plex"
	slackService "github.com/eschizoid/flixctl/slack/torrent"
	torrentService "github.com/eschizoid/flixctl/torrent"
	"github.com/spf13/cobra"
)

var DownloadTorrentCmd = &cobra.Command{
	Use:   "download",
	Short: "To Download A Torrent",
	Long:  "to download a torrent using Transmission client.",
	Run: func(cmd *cobra.Command, args []string) {
		envMagnetLink := os.Getenv("MAGNET_LINK")
		envDownloadDir := os.Getenv("DOWNLOAD_DIR")
		if envDownloadDir != "" && envMagnetLink == "" {
			// coming from web-hook
			downloadDir = envDownloadDir
			decodedEnvMagnetLink, err := base64.StdEncoding.DecodeString(envMagnetLink)
			if err != nil {
				fmt.Printf("Could not decode the magnet link: [%s]\n", err)
				panic(err)
			}
			magnetLink = string(decodedEnvMagnetLink)
		}
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := ec2.New(awsSession, awsSession.Config)
		instanceID := ec2Service.FetchInstanceID(svc, "plex")
		ec2status := ec2Service.Status(svc, instanceID)
		torrent := torrentService.TriggerDownload(magnetLink, downloadDir, ec2status)
		if torrent != nil {
			body, err := json.Marshal(torrent)
			if err != nil {
				fmt.Printf("Cannot parse response from transmission: [%s]\n", err)
				panic(err)
			}
			fmt.Println(string(body))
			if slackIncomingHookURL != "" {
				slackService.SendDownloadStart(*torrent.Name, slackIncomingHookURL)
			}
		} else {
			plexTorrentService.SendStatus("stopped", slackIncomingHookURL)
			fmt.Println("Plex status stopped")
		}
	},
}
