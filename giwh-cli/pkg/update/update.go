package update

import (
	"fmt"
	"github.com/fhluo/giwh/common/config"
	"github.com/fhluo/giwh/gacha-logs/gacha"
	"github.com/fhluo/giwh/gacha-logs/store"
	"github.com/fhluo/giwh/hyauth"
	"github.com/spf13/cobra"
	"log/slog"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update wish history",
		RunE: func(cmd *cobra.Command, args []string) error {
			var count int
			s := store.New(nil)
			err := s.LoadIfExists(config.WishHistoryPath.Get())
			if err != nil {
				return err
			}

			auths := hyauth.GenshinCN().Auths()
			if len(auths) == 0 {
				return fmt.Errorf("failed to find auth infos")
			}
			auth := auths[len(auths)-1]
			slog.Debug("update", "auth", auth)

			err = s.Update(auth, func(log gacha.Log) {
				count++
				fmt.Println(log.Name, log.UID, log.GachaType, log.Time)
			})
			if err != nil {
				return err
			}

			if count == 0 {
				fmt.Println("No items fetched. Your wish history is up to date.")
			} else {
				fmt.Printf("%d items fetched.\n", count)
			}

			return s.BackupAndSave(config.WishHistoryPath.Get())
		},
	}

	return cmd
}
