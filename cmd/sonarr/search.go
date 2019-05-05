package sonarr

import (
	"github.com/spf13/cobra"
)

var SearchSonarrCmd = &cobra.Command{
	Use:   "search",
	Short: "To Search Movies",
	Long:  "to search movies using sonarr client.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
