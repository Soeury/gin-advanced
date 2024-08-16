package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，用来保存所有的配置信息
var Conf = new(Config)

// 结构体保存配置文件信息
// 这里不管配置文件是什么类型 ，yaml , json , xml , binary ... ，结构体 tag 都使用 mapstructure
type Config struct {
	*StagingConfig `mapstructure:"staging"`
	*LogConfig     `mapstructure:"log"`
	*MysqlConfig   `mapstructure:"mysql"`
	*RedisConfig   `mapstructure:"redis"`
}

type StagingConfig struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

type LogConfig struct {
	Level       string `mapstructure:"level"`
	Filename    string `mapstructure:"filename"`
	Max_age     int    `mapstructure:"max_age"`
	Max_backups int    `mapstructure:"max_backups"`
	Max_size    int    `mapstructure:"max_size"`
}

type MysqlConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Dbname         string `mapstructure:"dbname"`
	Max_open_conns int    `mapstructure:"max_open_conns"`
	Max_idle_conns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Poolsize int    `mapstructure:"poolsize"`
}

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

	// 结构体保存配置文件信息需要 -> 将配置文件信息反序列化到全局变量 Conf 中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("found err in viper.Unmarshal : %v\n", err)
		return err
	}

	// 运行时实时读取配置文件的变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("the config file changed ...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("found err in viper.unmarshal : %v\n", err)
		}
	})
	return err
}
