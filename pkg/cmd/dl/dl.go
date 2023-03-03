package dl

import (
	"github.com/fhluo/giwh/pkg/cmd/dl/assets"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dl",
		Short: "Download assets",
	}

	cmd.AddCommand(assets.NewCmd())

	return cmd
}
