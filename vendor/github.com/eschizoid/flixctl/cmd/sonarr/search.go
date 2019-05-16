package sonarr

import (
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
		SearchShows()
	},
}

func SearchShows() {
	client, err := sonarr.New(os.Getenv("SONARR_URL"), os.Getenv("SONARR_API_KEY"))
	if err != nil {
		panic(err)
	}
	results, err := client.Search(keywords)
	if err != nil {
		panic(err)
	}
	slackService.SendShows(results, slackIncomingHookURL)
}
