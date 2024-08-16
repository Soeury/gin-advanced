package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// viper 配置文件

func main() {

	// 1. 读取配置文件 : 注意这里指定文件路径和文件名称只需要写一个就可以了
	viper.SetConfigFile("./config.yaml") // 指定文件路径
	// viper.SetConfigName("config.yaml")   // 配置文件名称

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // 类型断言
			fmt.Printf("config file not found : %s\n", err)
			return
		} else {
			fmt.Printf("exist another error : %s\n", err)
			return
		}
	}

	// 2. 写入配置文件 : 存储运行时的修改 (有safe前缀的都不会覆盖)
	viper.WriteConfig()             // 写入上面预定义的文件路径并覆盖
	viper.SafeWriteConfig()         // 写入上面预定义的文件路径但不覆盖
	viper.WriteConfigAs("path")     // 写入指定的路径并覆盖
	viper.SafeWriteConfigAs("path") // 写入指定的路径但不覆盖

	// 3. 支持在运行时实时读取配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config file changed : %s\n", e.Name)
	})

	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("version")) // 这里是 viper.GetString 不是 c.GetString
	})
	r.Run(":8080")

	// 4. 从 io.Reader 读取配置，这样会导致之前设置的 setconfigfile 失效，
	// 这样会从传递的 io.Reader 中读取配置，而不是之前设置的文件系统中读取配置(想直接使用文件)
	viper.SetConfigType("yaml")
	var example = []byte(`
				name: alice
				age:18
				eyes:blue
				hobbies:
				- swim
				- football
				- sleep
			`)

	viper.ReadConfig(bytes.NewBuffer(example))
	name := viper.GetString("name")
	fmt.Println(name)

	// 5. 注册和使用别名
	viper.RegisterAlias("apple", "melon") // 这里 apple 和 melon 建立了别名关系
	viper.Set("apple", 10)
	viper.Set("melon", 5)
	fmt.Printf("melon: %d\n", viper.GetInt("apple")) // melon = 5

	// 6. 支持使用环境变量
	viper.SetEnvPrefix("fff")  // 设置环境变量前缀
	viper.BindEnv("id")        // 绑定 id
	os.Setenv("fff_id", "111") // 程序内部设置环境变量 "fff_id" 的值为 "111"
	id := viper.GetString("id")
	fmt.Printf("id:%s\n", id)

	// 7. 使用 flags  将flag命令行参数绑定到viper中
	pflag.Int("flagname", 1234, "help msg for flagname") // 定义参数
	pflag.Parse()                                        // 解析参数
	viper.BindPFlags(pflag.CommandLine)                  // 参数绑定到viper
	value := viper.GetInt("filename")
	fmt.Printf("value:%d", value)

	// pflag也可以通过addgoflagset函数将 flag包中的参数绑定到viper中
	flag.Int("number", 666, "help msg for number")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	value = viper.GetInt("number")
	fmt.Printf("number:%d\n", value)

}
