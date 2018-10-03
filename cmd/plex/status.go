package plex

import (
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	"github.com/spf13/cobra"
)

var StatusPlexCmd = &cobra.Command{
	Use:   "status",
	Short: "To Get Plex Status",
	Long:  `to get the status of the Plex Media Center.`,
	Run: func(cmd *cobra.Command, args []string) {
		ec2Service.Status(Session, InstanceID)
	},
}
