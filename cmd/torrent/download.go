package torrent

import (
	"os"

	slackService "github.com/eschizoid/flixctl/slack/torrent"
	torrentService "github.com/eschizoid/flixctl/torrent"
	"github.com/spf13/cobra"
)

var DownloadTorrentCmd = &cobra.Command{
	Use:   "download",
	Short: "To Download a Torrent",
	Long:  "to download a torrent using Transmission client.",
	Run: func(cmd *cobra.Command, args []string) {
		envTorrentName := os.Getenv("TORRENT_NAME")
		envMagnetLink := os.Getenv("MAGNET_LINK")
		envDownloadDir := os.Getenv("DOWNLOAD_DIR")
		torrentService.TriggerDownload(envMagnetLink, argMagnetLink, envDownloadDir)
		if slackIncomingHookURL != "" {
			slackService.SendDownloadStart(envTorrentName, slackIncomingHookURL)
		}
	},
}
