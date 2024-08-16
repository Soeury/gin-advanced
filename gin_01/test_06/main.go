package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// sqlx库初识

var db *sqlx.DB

func initDB() (err error) {

	// 这里的 sqlx.Connect()相当于之前的   sql.Open() + db.Ping()
	// 也可以使用 MustConnect() 连接不成功就 panic
	// go语言里面很多调用都有 must 开头 , 不成功就 panic
	dsn := "root:123456@tcp(localhost:3306)/my_2?charset=utf8mb4&parseTime=True"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("found err in sqlx.connect : %v\n", err)
		return
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

func main() {

	err := initDB()
	if err != nil {
		fmt.Printf("found err in initDB : %v\n", err)
		return
	}
	fmt.Printf("connect my_2 successed!")
}
