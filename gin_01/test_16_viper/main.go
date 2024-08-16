package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Hugo struct {
}

func main() {

	// marshal     序列化:   [结构体，数组...]   ->   [json , xml , binary...]
	// unmarshal 反序列化:   [json , xml , binary...]  —>  [结构体，数组...]

	// 8. 远程 k/v 存储支持(etcd , consul)   -   加密 and 未加密
	var runtime_viper = viper.New()

	// 连接
	runtime_viper.AddRemoteProvider("etcd", "http://192.168.19.130:2379", "/config/hugo.yml")
	runtime_viper.SetConfigType("yml")

	// 读取远程配置
	err := runtime_viper.ReadRemoteConfig()
	if err != nil {
		fmt.Printf("unable to read remote config : %s\n", err)
		return
	}

	// 结构体接收配置文件
	var runtime_conf Hugo
	runtime_viper.Unmarshal(&runtime_conf)

	// 监控etcd配置文件的变化
	go func() {
		for {
			time.Sleep(time.Second * 5)

			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				fmt.Printf("unable to read remote config : %s\n", err)
				continue
			}

			// 配置文件更新时，自动加载到结构体
			runtime_viper.Unmarshal(&runtime_conf)
		}
	}()

	// 9. 从viper中获取值 , getxxx 方法在找不到值时会返回零值，可以用isSet(return bool)检查值是否存在
	viper.Set("number", 90)
	num := viper.GetInt("number")
	fmt.Printf("num:%d\n", num)

	if viper.IsSet("number") {
		fmt.Printf("number is set , value = %d\n", viper.GetInt("number"))
	} else {
		fmt.Println("number no exist")
	}

	// 10. 支持访问嵌套的值，每一级通过 . 的方式取出值，如果存在与分割的键路径匹配的键，则返回其值
	// 1. 设置路径
	viper.SetConfigFile("./config.json")

	// 2. 读取数据之后才能拿出数据
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // 类型断言
			fmt.Printf("config file not found : %s\n", err)
			return
		} else {
			fmt.Printf("exist another error : %s\n", err)
			return
		}
	}

	str := viper.GetString("datastore.metric.host") // 返回 "0.0.0.0"
	fmt.Printf("datastore.metric.host : %s\n", str)

	// 11. 从viper中提取子树
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // 类型断言
			fmt.Printf("config file not found : %s\n", err)
			return
		} else {
			fmt.Printf("exist another error : %s\n", err)
			return
		}
	}

	subv := viper.Sub("app.cache1")                                    // 此时 subv 代表了 [max-items: 100]  [item-size: 64] 这两个值
	fmt.Printf("app.cache.item-size : %d\n", subv.GetInt("item-size")) // subv.Get取出值

	// 12. viper 支持反序列化 将值解析到结构体或者map中
	viper.SetConfigFile("./config2.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // 类型断言
			fmt.Printf("config file not found : %s\n", err)
			return
		} else {
			fmt.Printf("exist another error : %s\n", err)
			return
		}
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Printf("found err in viper.unmarshal : %s\n", err)
		return
	} else {
		fmt.Printf("c: %+v\n", c)
	}

	// 13. 可以创建多个 viper 实例, viper是开箱即用的，即不需要配置或者初始化一个viper就可以使用
	x := viper.New()
	y := viper.New()

	x.SetDefault("number", 100)
	y.SetDefault("member", 99)
}

// *注意: 这里 tag 前面的值必须是 mapstructure
type Config struct {
	Port        int    `mapstructure:"port"`
	Version     string `mapstructure:"version"`
	MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	Dbname string `mapstructure:"dbname"`
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
}
