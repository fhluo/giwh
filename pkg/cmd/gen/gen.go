package gen

import (
	"github.com/fhluo/giwh/pkg/cmd/gen/lang"
	"github.com/fhluo/giwh/pkg/cmd/gen/menus"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate code",
	}

	cmd.AddCommand(lang.NewCmd())
	cmd.AddCommand(menus.NewCmd())

	return cmd
}