package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

// go - redis 连接
// 这里有一个问题就是: 连接不上

// 初始化一个变量
var rdb *redis.Client

// 连接redis数据库
func initClient() (err error) {

	// 设置相关参数
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.19.130:6379",
		Password: "123456",
		DB:       0,
		PoolSize: 100,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// 连接 redis 哨兵模式
// 连接 redis 集群

func main() {

	err := initClient()
	if err != nil {
		fmt.Printf("found err in initClient : %v\n", err)
		return
	}
	fmt.Println("connect with redis successed !")

	// 程序退出时释放相关资源，不能写在函数里面
	defer rdb.Close()
}
