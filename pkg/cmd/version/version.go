package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime/debug"
)

func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			if info, ok := debug.ReadBuildInfo(); !ok || info.Main.Version == "(devel)" {
				fmt.Println("dev")
			} else {
				fmt.Println("giwh version", info.Main.Version)
			}
		},
	}
	
	return cmd
}
