package radarr

import (
	"github.com/spf13/cobra"
)

var (
	RootRadarrCmd = &cobra.Command{
		Use:   "radarr",
		Short: "To Control Radarr",
	}
)

var (
	_ = func() struct{} {
		RootRadarrCmd.AddCommand(SearchRadarrCmd)
		return struct{}{}
	}()
)
