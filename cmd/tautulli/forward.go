package tautulli

import (
	"github.com/eschizoid/flixctl/tautulli"
	"github.com/spf13/cobra"
)

var ForwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "To Forward a Tautulli Event",
	Long:  "to forward a tautulli event to webhook implementation.",
	Run: func(cmd *cobra.Command, args []string) {
		tautulli.CreateEvent(args[1], args[2])
	},
}
