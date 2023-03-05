package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/cmd/giwh"
	"github.com/spf13/cobra"
	"runtime/debug"
)

var versionCmd = &cobra.Command{
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

func init() {
	main.rootCmd.AddCommand(versionCmd)
}
