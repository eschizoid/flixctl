package torrent

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var DownloadTorrentCmd = &cobra.Command{
	Use:   "download",
	Short: "To Download a Torrent",
	Long:  `to download a torrent using Transmission client.`,
	Run: func(cmd *cobra.Command, args []string) {
		transmission := exec.Command("transmission-remote", "--add", magnetLink)
		err := transmission.Start()
		if err != nil {
			fmt.Println("Could not download torrent using the given magnet link")
		}
	},
}
