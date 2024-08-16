package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

// go - redis 使用

// 初始化全局变量
var rdb *redis.Client

// 连接数据库
func initClient() (err error) {

	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.19.130",
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

func SetAndGet() {

	// get一个存在的值
	err := rdb.Set("score", 100, 0).Err()
	if err != nil {
		fmt.Printf("found err in rdb.Set : %v\n", err)
		return
	}

	val, err := rdb.Get("score").Result()
	if err != nil {
		fmt.Printf("found err in rdb.Get : %v\n", err)
		return
	}

	fmt.Println("score:", val)

	// get一个不存在的值
	val2, err := rdb.Get("name").Result()
	if err == redis.Nil {
		fmt.Println("not exist")
		return
	} else if err != nil {
		fmt.Printf("found err in rdb.get : %v\n", err)
		return
	} else {
		fmt.Println("name:", val2)
		return
	}
}

// redis - zset
// redis - pipeline
// redis - txpipeline
// redis - watch

func main() {

	err := initClient()
	if err != nil {
		fmt.Printf("found err in initClient : %v\n", err)
		return
	}
	defer rdb.Close()
	fmt.Printf("connect with redis successed!")

	SetAndGet()
}
