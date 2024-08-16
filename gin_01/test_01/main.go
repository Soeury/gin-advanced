package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin 源码解析 01

//  基数树(Radix Tree) 是一颗前缀树 , 我们注册路由的过程就是构造前缀树的过程, 具有公共前缀的节点会共享一个父节点,
//  路由器为每个请求方法管理一颗单独的树 , 每个树上的字节点都按照优先级排序
//  路由树是由一个一个节点组成的, gin框架中的路由树是由 node 结构体组成的

/*

	type node struct {
		path      string
		indices   string
		wildChild bool
		nType     nodeType
		priority  uint32
		children  []*node
		handlers  HandlersChain
		fullPath  string
	}

*/

func main() {

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})

	r.Run(":8080") //  Ctrl + 鼠标左键 点击

}
