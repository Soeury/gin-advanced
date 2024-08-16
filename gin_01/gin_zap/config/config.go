package config

import (
	"encoding/json"
)

// 定义一个全局配置
type Config struct {
	Mode       string `json:"mode"`
	Port       int    `json:"port"`
	*LogConfig `json:"log"`
}

// LogConfig 日志相关配置
type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxsize"`
	MaxAge     int    `json:"maxage"`
	MaxBackups int    `json:"maxbackups"`
}

var Conf = new(Config)

func Init() error {
	// 初始化LogConfig
	logConfig := &LogConfig{
		Level:      "info",
		Filename:   "D:\\M_GO\\GO_gin_advanced\\gin_01\\test_14\\log.txt",
		MaxSize:    1,
		MaxAge:     30,
		MaxBackups: 5,
	}

	// 初始化Config并设置LogConfig
	config := Config{
		Mode:      "debug", // debug , release , production
		Port:      8080,
		LogConfig: logConfig, // 确保LogConfig字段不是nil
	}

	// 将Config实例序列化为 JSON []byte
	json_byte, _ := json.Marshal(config)

	// 将json格式的[]byte切片解码到 conf 中
	return json.Unmarshal(json_byte, Conf)
}
