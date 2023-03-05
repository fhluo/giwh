package _import

import (
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/cmd/util"
	"github.com/fhluo/giwh/pkg/wish/pipeline"
	"github.com/fhluo/giwh/pkg/wish/repository"
	"github.com/spf13/cobra"
	"log"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import wish history",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filenames, err := util.ExpandPaths(args...)
			if err != nil {
				return err
			}

			items, err := repository.LoadItemsIfExits(config.WishHistoryPath)
			if err != nil {
				return err
			}

			p := pipeline.New(items)

			for _, filename := range filenames {
				items, err = repository.LoadItemsIfExits(filename)
				if err != nil {
					log.Fatalln(err)
				}
				p.Append(items...)
			}

			if err = repository.BackupAndSaveItems(config.WishHistoryPath, p.Items()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
