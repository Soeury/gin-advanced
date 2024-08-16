package main

import (
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// demo : gin中使用viper并且使用结构体保存变量信息

type Config struct {
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
}

var Conf = new(Config)

func main() {

	// 配置文件路径
	viper.SetConfigFile("./config.yaml")
	// 读取数据并检查
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("not found file err : %v\n", err)
			return
		}

		fmt.Printf("exist another err : %v\n", err)
		return
	}

	// 将读取的数据保存到全局变量 Conf 中
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("found err in unmarshal : %s", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()

	// 注册一个回调函数，以便在配置文件发生修改时调用，特别适合需要更新配置而不需要重启服务时调用
	// 配置文件变化后要同步到Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("the configed changed...")
		if err := viper.Unmarshal(&Conf); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err : %s", err))
		}
	})

	// version 的返回值会随着配置文件的变化而变化
	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, Conf.Version)
	})

	if err := r.Run(fmt.Sprintf(":%d", Conf.Port)); err != nil {
		panic(err)
	}
}
