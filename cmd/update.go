package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sort"
	"strconv"
)

func Default() api.Region {
	if _, err := os.Stat(api.CN.OutputLogPath); err == nil {
		return api.CN
	}

	if _, err := os.Stat(api.OS.OutputLogPath); err == nil {
		return api.OS
	}

	return api.OS
}

func FetchAllWishHistory(ctx *api.Context, items wish.Items) (wish.Items, error) {
	visit := make(map[int64]bool)
	for _, item := range items {
		visit[item.ID()] = true
	}
	sort.Sort(sort.Reverse(items))

	for _, type_ := range wish.SharedTypes {
		fmt.Printf("Fetching the wish history of %s.\n", type_.GetSharedWishName())

		x := items.FilterByWishType(type_)
		if len(x) != 0 {
			result, err := ctx.WishType(strconv.Itoa(type_)).Size(10).Begin(x[0].Item.ID).FetchAll()
			if err != nil {
				return nil, err
			}

			items = append(lo.Map(lo.Reverse(result), func(item *api.Item, _ int) wish.Item {
				return wish.Item{Item: item}
			}), items...)
		} else {
			result, err := ctx.WishType(strconv.Itoa(type_)).Size(10).End("0").FetchAll()
			if err != nil {
				return nil, err
			}

			items = append(items, lo.Map(result, func(item *api.Item, _ int) wish.Item {
				return wish.Item{Item: item}
			})...)
		}

	}

	return items, nil
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update wish history",
	Run: func(cmd *cobra.Command, args []string) {
		base, err := Default().GetAPIBase()
		if err != nil {
			log.Fatalln(err)
		}

		ctx, err := api.New(base)
		if err != nil {
			log.Fatalln(err)
		}

		uid, err := ctx.GetUID()
		if err != nil {
			log.Fatalln(err)
		}

		items := config.WishHistory.FilterByUID(uid)
		result, err := FetchAllWishHistory(ctx, items)
		if err != nil {
			log.Fatalln(err)
		}

		count := len(result) - len(items)
		if count == 0 {
			fmt.Println("No items fetched. Your wish history is up to date.")
			return
		}

		config.WishHistory = append(config.WishHistory, result...)
		if err := config.SaveWishHistory(); err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%d items fetched.\n", count)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
