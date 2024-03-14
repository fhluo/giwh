package dl

import (
	"github.com/spf13/cobra"
	"tools/dl/assets"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dl",
		Short: "Download assets",
	}

	cmd.AddCommand(assets.NewCmd())

	return cmd
}
