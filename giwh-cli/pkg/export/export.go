package export

import (
	"fmt"
	"github.com/fhluo/giwh/common/config"
	"github.com/fhluo/giwh/gacha-logs/pipeline"
	"github.com/fhluo/giwh/gacha-logs/store"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log"
	"strconv"
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
			s := store.New(nil)

			err := s.LoadIfExists(config.WishHistoryPath.Get())
			if err != nil {
				return err
			}

			p := pipeline.New(s.Unique())

			if cmd.Flags().Changed("uid") {
				p = p.FilterByUID(strconv.Itoa(uid))
			}

			if cmd.Flags().Changed("wish") {
				p = p.FilterBySharedWish(lo.Map(wishes, func(i int, _ int) string {
					return strconv.Itoa(i)
				})...)
			}

			if p.Len() == 0 {
				log.Fatalln("No such wish history.")
			}

			if err = store.New(p.Logs()).Save(args[0]); err != nil {
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
