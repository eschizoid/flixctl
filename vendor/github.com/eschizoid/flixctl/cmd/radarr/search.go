package radarr

import (
	"github.com/spf13/cobra"
)

var SearchRadarrCmd = &cobra.Command{
	Use:   "search",
	Short: "To Search Shows",
	Long:  "to search shows using radarr client.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
