package plex

import (
	"github.com/spf13/cobra"
)

var RootPlexCmd = &cobra.Command{
	Use:   "plex",
	Short: "To Control Plex Media Center",
}

func init() {
	RootPlexCmd.AddCommand(StartPlexCmd, StopPlexCmd, StatusPlexCmd)
}
