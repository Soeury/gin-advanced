package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	// os.Args[0] 默认是程序的名字
	// 如果执行以下代码 :
	// go run "d:\M_GO\GO_gin_advanced\gin_01\test_22_flag\main.go" 1 2 3 4 5
	// 会打印下面6条消息
	// args[0] = project_name
	// args[1] = 1
	// args[2] = 2
	// args[3] = 3
	// args[4] = 4
	// args[5] = 5

	fmt.Println(os.Args[0])
	if len(os.Args) > 0 {
		for i, v := range os.Args {
			fmt.Printf("args[%d]=%v\n", i, v)
		}
	}

	// 除了可以使用 os.Args , 还可以使用 flag 和 pflag 包来获取命令行参数
	// 1. flag
	// value := flag.Type("name" , "default_value" , "message") (*type)
	// flag.Parse()
	// 最后可以通过 *value 的方式取出值 ， 命令行没有传入值就采用默认值
	name1 := flag.String("name", "alice", "somebody's characteristic")
	flag.Parse() // 好像是只能 Parse 一次
	fmt.Printf("%s\n", *name1)

	var age1 int
	flag.IntVar(&age1, "age", 18, "somebody's lifelong")
	flag.Parse()
	fmt.Printf("%d\n", age1)

	fmt.Println("-----------------------------------------")

	// A Demo
	// 定义命令行的变量
	var name string
	var age int
	var married bool
	var delay time.Duration
	flag.StringVar(&name, "name", "peter", "somebody's characteristic")
	flag.IntVar(&age, "age", 18, "somebody's lifelong")
	flag.BoolVar(&married, "married", false, "somebody's kiss or not")
	flag.DurationVar(&delay, "delay", 0, "the delay time gap")

	// 解析参数
	flag.Parse()

	// go run "d:\M_GO\GO_gin_advanced\gin_01\test_22_flag\main.go" -h
	// 上面这段代码会将设置的变量信息全部打印出来

	// * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
	// 在执行代码的后面加上 --name=value 相当于传参，这时候不会采用默认值

	fmt.Println(name, age, married, delay)
	// 返回命令行参数后的其他参数
	fmt.Println(flag.Args())
	// 返回命令行参数后的其他参数的个数
	fmt.Println(flag.NArg())
	// 返回使用的命令行参数的个数
	fmt.Println(flag.NFlag())

}
