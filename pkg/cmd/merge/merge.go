package merge

import (
	"github.com/fhluo/giwh/pkg/cmd/util"
	"github.com/fhluo/giwh/pkg/wish/pipeline"
	"github.com/fhluo/giwh/pkg/wish/repository"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"log"
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

			p := pipeline.New(nil)

			for _, filename := range filenames {
				items, err := repository.LoadItemsIfExits(filename)
				if err != nil {
					return err
				}
				p.Append(items...)
			}

			p.SortByIDDescending()
			if err = repository.SaveItems(outputFilename, p.Unique().Items()); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFilename, "output", "o", "", "specify output filename")
	if err := cmd.MarkFlagRequired("output"); err != nil {
		slog.Error(err.Error(), nil)
		os.Exit(1)
	}

	return cmd
}
