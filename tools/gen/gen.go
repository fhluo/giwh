package gen

import (
	"github.com/spf13/cobra"
	"tools/gen/lang"
	"tools/gen/locales"
	"tools/gen/menus"
	"tools/gen/wishes"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate code",
	}

	cmd.AddCommand(lang.NewCmd())
	cmd.AddCommand(menus.NewCmd())
	cmd.AddCommand(locales.NewCmd())
	cmd.AddCommand(wishes.NewCmd())

	return cmd
}
