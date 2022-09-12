package cmd

import (
	"fmt"
	"github.com/fhluo/giwh/internal/config"
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/samber/lo"
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

func FetchAllWishHistory(ctx *api.Context, p pipeline.Pipeline) (pipeline.Pipeline, error) {
	visit := make(map[int64]bool)
	for _, item := range p.Elements() {
		visit[item.ID] = true
	}
	p = p.SortedByIDDescending()

	for _, type_ := range []string{api.CharacterEventWish, api.WeaponEventWish, api.StandardWish, api.BeginnersWish} {
		fmt.Printf("Fetching the wish history of %s.\n", type_)

		x := p.FilterByWishType(type_).Items()
		if len(x) != 0 {
			result, err := ctx.WishType(type_).Size(10).Begin(x[0].ID).FetchAll()
			if err != nil {
				return p, err
			}

			r, err := pipeline.ItemsTo(lo.Reverse(result), pipeline.NewElement)
			if err != nil {
				return p, err
			}

			p = p.Append(r)
		} else {
			result, err := ctx.WishType(type_).Size(10).End("0").FetchAll()
			if err != nil {
				return p, err
			}

			r, err := pipeline.ItemsTo(lo.Reverse(result), pipeline.NewElement)
			if err != nil {
				return p, err
			}

			p = p.Append(r)
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

		count := len(result.Elements()) - len(config.WishHistory.Elements())
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
