package wishes

import (
	_ "embed"
	"fmt"
	"github.com/dop251/goja"
	"github.com/fhluo/giwh/pkg/wish"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
	"os"
)

var (
	//go:embed wishes.js
	wishesJS string

	Wishes       map[string][]wish.Type
	SharedWishes map[string][]wish.Type

	log = slog.With("gen wishes")
)

func init() {
	vm := goja.New()

	_, err := vm.RunString(wishesJS)
	if err != nil {
		log.Error(err.Error(), nil)
		os.Exit(1)
	}

	err = mapstructure.WeakDecode(vm.Get("itemTypeMap").Export(), &SharedWishes)
	if err != nil {
		log.Error(err.Error(), nil)
		os.Exit(1)
	}

	err = mapstructure.WeakDecode(vm.Get("itemTypeNameMap").Export(), &Wishes)
	if err != nil {
		log.Error(err.Error(), nil)
		os.Exit(1)
	}

	for key := range Wishes {
		slices.SortFunc(Wishes[key], func(a, b wish.Type) bool {
			return a.Key < b.Key
		})
	}

	for key := range SharedWishes {
		slices.SortFunc(SharedWishes[key], func(a, b wish.Type) bool {
			return a.Key < b.Key
		})
	}
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wishes",
		Short: "Generate pkg/wishes/types.go",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(Wishes)
			return nil
		},
	}

	return cmd
}
