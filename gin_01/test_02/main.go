package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// gin 源码 - 中间件
//    ---  c.Next() 调用接下来的所有函数
//    ---  c.Abort() 不调用接下来的函数

//    ---  c.Set("key" , "value") 设置值
//    ---  value , ok := c.Get("key") 取出值

func temp1(c *gin.Context) {

	fmt.Println("this is temp1")
}

func temp2(c *gin.Context) {

	fmt.Println("this is temp2 before")
	c.Next()
	fmt.Println("this is temp2 after")
}

func temp3(c *gin.Context) {

	fmt.Println("this is temp3")
}

func temp4(c *gin.Context) {

	fmt.Println("this is temp4")
	c.Set("name", "chen")
}

func temp5(c *gin.Context) {

	fmt.Println("this is temp5")
	v, ok := c.Get("name")
	if !ok {
		fmt.Printf("found error in c.Get")
		return
	}

	fmt.Println("name : ", v)
}

func main() {

	r := gin.Default()

	shopGroup := r.Group("/shop", temp1, temp2)
	shopGroup.Use(temp3)
	{
		shopGroup.GET("/buy", temp4, temp5)
	}

	r.Run(":8080")

}
