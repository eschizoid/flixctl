package sonarr

import (
	"github.com/spf13/cobra"
)

var (
	RootSonarrCmd = &cobra.Command{
		Use:   "sonarr",
		Short: "To Control Sonarr",
	}
)

var (
	_ = func() struct{} {
		RootSonarrCmd.AddCommand(SearchSonarrCmd)
		return struct{}{}
	}()
)
