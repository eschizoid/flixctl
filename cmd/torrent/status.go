package torrent

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var StatusTorrentCmd = &cobra.Command{
	Use:   "status",
	Short: "To Show Torrents Status",
	Long:  `to show the status of the torrents being downloaded`,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("transmission-remote", "--torrent=active", "--list").Output()
		if err != nil {
			fmt.Println("Could not list torrents being downloaded")
		}
		fmt.Println(string(out))
	},
}
