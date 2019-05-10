package radarr

import (
	"os"

	slackService "github.com/eschizoid/flixctl/slack/radarr"
	"github.com/jrudio/go-radarr-client"
	"github.com/spf13/cobra"
)

var SearchRadarrCmd = &cobra.Command{
	Use:   "search",
	Short: "To Search Movies",
	Long:  "to search movies using radarr client.",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := radarr.New(os.Getenv("RADARR_URL"), os.Getenv("RADARR_API_KEY"))
		if err != nil {
			panic(err)
		}
		results, err := client.Search(keywords)
		if err != nil {
			panic(err)
		}
		slackService.SendMovies(results, slackIncomingHookURL)
	},
}
