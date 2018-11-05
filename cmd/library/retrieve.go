package library

import (
	"github.com/spf13/cobra"
)

var RetrieveLibraryCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "To Retrieve a File From Media Library",
	Long:  "to retrieve a movie from media library.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
