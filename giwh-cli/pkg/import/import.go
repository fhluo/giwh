package _import

import (
	"github.com/fhluo/giwh/common/config"
	"github.com/fhluo/giwh/gacha-logs/store"
	"github.com/fhluo/giwh/giwh-cli/pkg/util"
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

			s := store.New(nil)
			err = s.LoadIfExists(config.WishHistoryPath.Get())
			if err != nil {
				return err
			}

			for _, filename := range filenames {
				err = s.LoadIfExists(filename)
				if err != nil {
					log.Fatalln(err)
				}
			}

			if err = s.BackupAndSave(config.WishHistoryPath.Get()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
