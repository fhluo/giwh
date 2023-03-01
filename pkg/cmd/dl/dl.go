package dl

import (
	"github.com/fhluo/giwh/pkg/cmd/dl/icons"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dl",
		Short: "Download assets",
	}

	cmd.AddCommand(icons.NewCmd())

	return cmd
}
