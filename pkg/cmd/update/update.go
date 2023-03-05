package update

import (
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/fhluo/giwh/pkg/wish/repository"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update wish history",
		RunE: func(cmd *cobra.Command, args []string) error {
			var count int
			_, err := repository.UpdateItems(config.WishHistoryPath, func(item wish.Item) {
				count++
				fmt.Println(item.Name, item.UID, item.WishType, item.Time)
			})
			if err != nil {
				return err
			}

			if count == 0 {
				fmt.Println("No items fetched. Your wish history is up to date.")
			} else {
				fmt.Printf("%d items fetched.\n", count)
			}

			return nil
		},
	}

	return cmd
}
