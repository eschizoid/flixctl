package plex

import (
	"github.com/aws/aws-sdk-go/aws/session"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	"github.com/spf13/cobra"
)

var RootPlexCmd = &cobra.Command{
	Use:   "plex",
	Short: "To Control Plex Media Center",
}

var Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var InstanceID = ec2Service.FetchInstanceID(Session, "plex")

func init() {
	RootPlexCmd.AddCommand(StartPlexCmd, StopPlexCmd, StatusPlexCmd)
}
