package ombi

import (
	"github.com/spf13/cobra"
)

var RequestOmbiCmd = &cobra.Command{
	Use:   "request",
	Short: "To Request Movies or Shows",
	Long:  "to request movies or shows via ombi.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
