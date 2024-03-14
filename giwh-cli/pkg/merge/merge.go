package merge

import (
	"github.com/fhluo/giwh/gacha-logs/store"
	"github.com/fhluo/giwh/giwh-cli/pkg/util"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
)

func NewCmd() *cobra.Command {
	var outputFilename string

	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge wish histories",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			filenames, err := util.ExpandPaths(args...)
			if err != nil {
				log.Fatalln(err)
			}

			s := store.New(nil)

			for _, filename := range filenames {
				if err := s.LoadIfExists(filename); err != nil {
					return err
				}
			}

			if err = s.Save(outputFilename); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFilename, "output", "o", "", "specify output filename")
	if err := cmd.MarkFlagRequired("output"); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return cmd
}
