package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

//  ----  go - 默认 logger 日志库
// -1 优势: 简便
// -2 劣势: 仅限基本的日志级别，对于错误日志的处理比较不方便，不提供日志切割的功能

// 设置 go 内置的 logger 日志记录器
func SetUpLogger() {

	path := "./log.txt"
	logFileLoc, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	if err != nil {
		fmt.Printf("found err in os.openfile : %v\n", err)
		return
	}
	log.SetOutput(logFileLoc)
}

// 日志演示 - 只需要在打印输出前面加上 log. 即可
func logDemo(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error in url %s : %s", url, err.Error())
	} else {
		log.Printf("status code for %s : %s", url, resp.Status)
		resp.Body.Close()
	}
}

func main() {

	SetUpLogger()
	logDemo("https://www.baidu.com")
	logDemo("www.google.com")
}
