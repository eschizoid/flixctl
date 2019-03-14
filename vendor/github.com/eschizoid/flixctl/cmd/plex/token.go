package plex

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jrudio/go-plex-client" //nolint:goimports
	"github.com/spf13/cobra"
)

var TokenPlexCmd = &cobra.Command{
	Use:   "token",
	Short: "To Get Plex Token",
	Long:  "to get a Plex token for API calls.",
	Run: func(cmd *cobra.Command, args []string) {
		plexClient, err := plex.SignIn(os.Getenv("PLEX_USER"), os.Getenv("PLEX_PASSWORD"))
		ShowError(err)
		m := make(map[string]interface{})
		m["plex_token"] = plexClient.Token
		jsonString, _ := json.Marshal(m)
		fmt.Println(string(jsonString))
	},
}
