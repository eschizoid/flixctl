package nzb

import (
	"github.com/spf13/cobra"
)

var StatusNzbCmd = &cobra.Command{
	Use:   "status",
	Short: "To Show NZB Status",
	Long:  "to show the status of the nzb files being downloaded.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
