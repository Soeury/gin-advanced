package main

import "github.com/spf13/viper"

// 关于 viper 的几个注意事项
type Config struct {
	Name    string `mapstructure:"name"`
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
}

func main() {

	// 1. 使用相对/绝对路径添加配置文件
	viper.SetConfigFile("relative_path")
	viper.SetConfigFile("absolutely_path")

	// 2. 指定配置文件名称和位置 ， viper 自行查找文件
	viper.SetConfigName("file_name")
	viper.AddConfigPath("path")

	// 3. 下面这种方式配合远程中心使用，告诉viper 要解析的是哪种类型的配置文件
	viper.SetConfigType("type")

	// 4. 通过结构体保存配置文件信息需要先将配置文件信息进行反序列化
	var p Config
	viper.Unmarshal(p)

	// 5. 使用结构体 tag 前面的类型固定为 mapstructure
}
