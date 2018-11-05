package library

import (
	"github.com/spf13/cobra"
)

var SearchLibraryCmd = &cobra.Command{
	Use:   "download",
	Short: "To Search For A File In The Media Library",
	Long:  "to search for a movie or show in media library.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
