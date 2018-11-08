package library

import (
	"github.com/spf13/cobra"
)

var RootLibraryCmd = &cobra.Command{
	Use:   "library",
	Short: "To Control Media Library",
}

func init() {
	RootLibraryCmd.AddCommand(RetrieveLibraryCmd, ArchiveLibraryCmd, SearchLibraryCmd, RefreshLibraryCmd)
}
