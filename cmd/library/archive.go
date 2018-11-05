package library

import (
	"github.com/spf13/cobra"
)

var ArchiveLibraryCmd = &cobra.Command{
	Use:   "archive",
	Short: "To Archive a File To Media Library",
	Long:  "to archive a movie or show to media library.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
