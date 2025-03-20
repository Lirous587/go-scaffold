package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func initViper() error {
	viper.SetConfigFile("manifest/config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	//实时监控配置文件
	viper.WatchConfig()

	//配置文件修改之后的回调函数
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件已修改: %s (%s)\n", e.Name, e.Op)

		// 热重载配置
		if err := loadConfig(); err != nil {
			fmt.Printf("重新加载配置失败: %v\n", err)
		}
	})
	return nil
}
