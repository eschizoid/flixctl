package ombi

import (
	"github.com/spf13/cobra"
)

var (
	RootOmbiCmd = &cobra.Command{
		Use:   "ombi",
		Short: "To Control Ombi",
	}
)

var (
	_ = func() struct{} {
		RootOmbiCmd.AddCommand(RequestOmbiCmd)
		return struct{}{}
	}()
)
