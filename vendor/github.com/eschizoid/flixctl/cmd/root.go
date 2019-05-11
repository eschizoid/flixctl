package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eschizoid/flixctl/cmd/admin"
	"github.com/eschizoid/flixctl/cmd/library"
	"github.com/eschizoid/flixctl/cmd/nzb"
	"github.com/eschizoid/flixctl/cmd/ombi"
	"github.com/eschizoid/flixctl/cmd/plex"
	"github.com/eschizoid/flixctl/cmd/radarr"
	"github.com/eschizoid/flixctl/cmd/sonarr"
	"github.com/eschizoid/flixctl/cmd/torrent"
	"github.com/spf13/cobra"
)

var FlixctlCmd = &cobra.Command{
	Use: "flixctl",
	Long: `To Control The Following flixctl Components:
  * Library
  * NZB Client
  * Ombi
  * Plex
  * Radarr
  * Sonarr
  * Torrent Client`,
}

var (
	// BUILD and VERSION are set during build
	BUILD   string
	VERSION string
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "To Get flixctl Version",
	Long:  "to get flixctl version",
	Run: func(cmd *cobra.Command, args []string) {
		m := make(map[string]string)
		m["version"] = VERSION
		m["build"] = BUILD
		jsonString, _ := json.Marshal(m)
		fmt.Println(string(jsonString))
	},
}

var (
	_ = func() struct{} {
		FlixctlCmd.AddCommand(
			VersionCmd,
			admin.FlixctlAdminRootCmd,
			library.RootLibraryCmd,
			nzb.RootNzbCmd,
			ombi.RootOmbiCmd,
			plex.RootPlexCmd,
			radarr.RootRadarrCmd,
			sonarr.RootSonarrCmd,
			torrent.RootTorrentCmd,
		)
		return struct{}{}
	}()
)

func Execute(version string, build string) {
	VERSION = version
	BUILD = build

	if err := FlixctlCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
