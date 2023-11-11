package config

import (
	"errors"
	"github.com/fhluo/giwh/i18n"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path/filepath"
)

var (
	Dir = filepath.Join(os.Getenv("LOCALAPPDATA"), "giwh")

	Path = filepath.Join(Dir, "config.toml")

	WishHistoryPath = NewItem("wish_history_path", filepath.Join(Dir, "wish_history.json"))
	DBPath          = NewItem("db_path", filepath.Join(Dir, "wish_history.json"))
	Language        = NewItem("language", i18n.Default().Tag().String())
)

func init() {
	// 创建配置文件目录
	err := os.MkdirAll(Dir, 0666)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// 添加配置文件目录，设置配置文件名
	viper.AddConfigPath(Dir)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// 读取配置文件
	err = viper.ReadInConfig()
	if err == nil {
		return
	}

	// 读取配置文件失败，判断错误类型
	var notFoundErr viper.ConfigFileNotFoundError
	if !errors.As(err, &notFoundErr) {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// 未找到配置文件，创建配置文件
	slog.Info(notFoundErr.Error())
	if err = viper.WriteConfigAs(Path); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

// Save 保存配置文件
func Save() {
	// 写入配置
	err := viper.WriteConfig()
	if err == nil {
		return
	}

	// 写入配置失败，判断错误类型
	var notFoundErr viper.ConfigFileNotFoundError
	if !errors.As(err, &notFoundErr) {
		slog.Warn(err.Error())
		return
	}

	// 重新尝试写入配置
	if err = viper.WriteConfigAs(Path); err != nil {
		slog.Warn(err.Error())
	}
}
