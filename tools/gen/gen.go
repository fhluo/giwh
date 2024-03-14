package gen

import (
	"github.com/fhluo/giwh/giwh-cli/pkg/gen/lang"
	"github.com/fhluo/giwh/giwh-cli/pkg/gen/locales"
	"github.com/fhluo/giwh/giwh-cli/pkg/gen/menus"
	"github.com/fhluo/giwh/giwh-cli/pkg/gen/wishes"
	"github.com/spf13/cobra"
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
