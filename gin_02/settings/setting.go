package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func ConfigInit() (err error) {

	// 指定文件路径 + 读取配置信息 + 错误检查
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("not found config file : %s\n", err)
			return err
		}
		fmt.Printf("exist another err : %s\n", err)
		return err
	}

	// 运行时实时读取配置文件的变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("the config file changed ...")
	})
	return err
}
