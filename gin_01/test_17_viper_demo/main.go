package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// demo : gin框架搭建的web项目中使用viper

func main() {

	// 指定配置文件路径
	viper.SetConfigFile("./config.yaml")
	// 读取配置信息并检验
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("not found file : %v\n", err)
			return
		}
		fmt.Printf("exist another err : %v\n", err)
		return
	}

	// 监控配置文件变化
	viper.WatchConfig()

	// 访问version的返回值会随着配置文件的变化而变化
	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("version"))
	})

	if err := r.Run(fmt.Sprintf(":%d", viper.GetInt("port"))); err != nil {
		panic(err)
	}
}
