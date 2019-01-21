package nzb

import (
	"github.com/spf13/cobra"
)

var (
	RootNzbCmd = &cobra.Command{
		Use:   "nzb",
		Short: "To Control Nzb Client",
	}
)

var (
	_ = func() struct{} {
		RootNzbCmd.AddCommand(StatusNzbCmd)
		return struct{}{}
	}()
)