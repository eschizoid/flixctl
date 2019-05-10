package sonarr

import (
	"fmt"
	"os"

	slackService "github.com/eschizoid/flixctl/slack/sonarr"
	"github.com/jrudio/go-sonarr-client"
	"github.com/spf13/cobra"
)

var SearchSonarrCmd = &cobra.Command{
	Use:   "search",
	Short: "To Search Shows",
	Long:  "to search shows using sonarr client.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := sonarr.New(fmt.Sprintf("htttp://%s:%d/sonarr", os.Getenv("FLIXCTL_HOST"), 9443), os.Getenv("SONARR_API_KEY"))
		if err != nil {
			panic(err)
		}
		results, err := client.Search(keywords)
		if err != nil {
			panic(err)
		}
		slackService.SendShows(results, slackIncomingHookURL)
	},
}
