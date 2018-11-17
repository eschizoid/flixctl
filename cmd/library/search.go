package library

import (
	"fmt"

	"github.com/eschizoid/flixctl/library"
	"github.com/spf13/cobra"
)

var FindLibraryCmd = &cobra.Command{
	Use:   "initiate",
	Short: "To Find A Movie Or Show",
	Long:  "to find a movie or show in the library.",
	Run: func(cmd *cobra.Command, args []string) {
		glacierResults := library.FindInLibrary(query)
		fmt.Println(glacierResults)
		plexResults := library.FindInPlex(query)
		fmt.Println(plexResults)
	},
}
