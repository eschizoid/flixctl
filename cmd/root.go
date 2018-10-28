package cmd

import (
	"fmt"
	"os"

	"github.com/eschizoid/flixctl/cmd/plex"
	"github.com/eschizoid/flixctl/cmd/torrent"
	"github.com/spf13/cobra"
)

var FlixctlCmd = &cobra.Command{
	Use: "flixctl",
	Long: `To Control The Following flixctl Components:
  * Plex
  * Tautulli
  * Torrent`,
}

func init() {
	FlixctlCmd.AddCommand(plex.RootPlexCmd, torrent.RootTorrentCmd)
}

func Execute() {
	if err := FlixctlCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
