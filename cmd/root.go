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

var (
	// BUILD and VERSION are set during build
	BUILD   string
	VERSION string
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "To Get flixctl version",
	Long:  "to get flixctl version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(" " + VERSION)
		fmt.Println(" " + BUILD)
	},
}

func init() {
	FlixctlCmd.AddCommand(VersionCmd, plex.RootPlexCmd, torrent.RootTorrentCmd)
}

func Execute(version string, build string) {
	VERSION = version
	BUILD = build

	if err := FlixctlCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
