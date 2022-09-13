package cmd

import (
	"github.com/fhluo/giwh/pkg/api"
	"github.com/fhluo/giwh/pkg/pipeline"
	"github.com/fhluo/giwh/pkg/repository"
	"github.com/spf13/cobra"
	"log"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge wish histories",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filenames, err := ExpandPaths(args...)
		if err != nil {
			log.Fatalln(err)
		}

		var result []*api.Item

		for _, filename := range filenames {
			items, err := repository.Load(filename)
			if err != nil {
				log.Fatalln(err)
			}
			result = append(result, items...)
		}

		p, err := pipeline.New(result)
		if err != nil {
			log.Fatalln(err)
		}

		if err = repository.Save(outputFilename, p.Unique().SortedByIDDescending().Items()); err != nil {
			log.Fatalln(err)
		}
	},
}

var outputFilename string

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().StringVarP(&outputFilename, "output", "o", "", "specify output filename")
	if err := mergeCmd.MarkFlagRequired("output"); err != nil {
		log.Fatalln(err)
	}
}
