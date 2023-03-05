package export

import (
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/fhluo/giwh/pkg/wish/pipeline"
	"github.com/fhluo/giwh/pkg/wish/repository"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log"
)

func NewCmd() *cobra.Command {
	var (
		uid    int
		wishes []int
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export wish history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := repository.LoadItems(config.WishHistoryPath)
			if err != nil {
				return err
			}

			p := pipeline.New(items)

			if cmd.Flags().Changed("uid") {
				p = p.FilterByUID(uid)
			}

			if cmd.Flags().Changed("wish") {
				p = p.FilterBySharedWish(lo.Map(wishes, func(i int, _ int) wish.Type {
					return wish.Type(i)
				})...)
			}

			if p.Len() == 0 {
				log.Fatalln("No such wish history.")
			}

			if err = repository.SaveItems(args[0], p.Items()); err != nil {
				return err
			} else {
				fmt.Printf("%d items exported.\n", p.Len())
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&uid, "uid", "u", 0, "specify uid")
	cmd.Flags().IntSliceVarP(&wishes, "wishes", "w", nil, "specify wish types")

	return cmd
}
