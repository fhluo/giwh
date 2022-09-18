package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/spf13/cobra"
	"log"
	"os"
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

func FetchAllWishHistory(ctx *api.Context, p pipeline.Items) (pipeline.Items, error) {
	visit := make(map[int64]bool)
	for _, item := range p {
		visit[item.ID] = true
	}
	p.SortByIDDescending()

	for _, wishType := range api.SharedWishTypes {
		fmt.Printf("Fetching the wish history of %s.\n", wishType)

		x := p.FilterBySharedWishType(wishType)
		if len(x) != 0 {
			result, err := ctx.WishType(wishType).Size(10).Begin(x[0].ID).FetchAll()
			if err != nil {
				return p, err
			}

			p = p.Append(result...)
		} else {
			result, err := ctx.WishType(wishType).Size(10).End(0).FetchAll()
			if err != nil {
				return p, err
			}

			p = p.Append(result...)
		}

	}

	return p, nil
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

		result, err := FetchAllWishHistory(ctx, config.WishHistory.FilterByUID(uid))
		if err != nil {
			log.Fatalln(err)
		}

		count := len(result) - len(config.WishHistory)
		if count == 0 {
			fmt.Println("No items fetched. Your wish history is up to date.")
			return
		}

		config.WishHistory = result
		if err := config.SaveWishHistory(); err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%d items fetched.\n", count)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
