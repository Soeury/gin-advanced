package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// mysql 预处理
//   ---- 批量处理同一个sql语句的情况下采用预处理能够提高处理速率
//   ---- 避免sql注入
//   *---- 不要自己拼接 sql语句

var db *sql.DB

// 连接数据库 - 记得加上驱动
func initMysql() (err error) {

	dsn := "root:123456@tcp(localhost:3306)/my_2"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("found err in db.ping : %v\n", err)
		return
	}
	return
}

type Users struct {
	Id   int
	Name string
	Age  int
}

// 预处理数据查询
func pre_query() {

	sqlStr := "select id , name , age from users where id > ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("found err in db.prepare : %v\n", err)
		return
	}
	defer stmt.Close() // 释放连接

	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("found err in stmt.query : %v\n", err)
		return
	}
	defer rows.Close() // 释放连接

	// 循环读取数据
	for rows.Next() {
		var u Users
		err := rows.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			fmt.Printf("found err in rows.scan : %v\n", err)
			return
		}
		fmt.Printf("id:%d\tname:%s\tage:%d\n", u.Id, u.Name, u.Age)
	}
}

// 预处理插入
func pre_insert() {

	sqlStr := "insert into users(name , age) values(? , ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("found err in db.exec : %v\n", err)
		return
	}
	defer stmt.Close()

	// 插入两条数据
	_, err = stmt.Exec("jun", 13)
	if err != nil {
		fmt.Printf("found err in stmt.exec : %v\n", err)
		return
	}

	_, err = stmt.Exec("ong", 11)
	if err != nil {
		fmt.Printf("found err in stmt.exec : %v\n", err)
		return
	}
	fmt.Println("insert sucdessed!")
}

// sql注入示范
func sqlInject(name string) {

	sqlStr := fmt.Sprintf("select id , name , age from users where name='%s'", name)
	fmt.Printf("sql:%s\n", sqlStr)
	var u Users
	err := db.QueryRow(sqlStr).Scan(&u.Id, &u.Name, &u.Age)
	if err != nil {
		fmt.Printf("found err in db.queryrow : %v\n", err)
		return
	}
	fmt.Printf("%#v\n", u)
}

func main() {

	err := initMysql()
	if err != nil {
		fmt.Printf("found err in initmysql : %v\n", err)
		return
	}

	defer db.Close()
	fmt.Println("connect with my_2 successed!")

	// 查询
	pre_query()

	// 插入
	pre_insert()

	// sql注入
	sqlInject("xxx' or 1=1#")
	//sqlInject("xxx' union select * from users #")
	//sqlInject("xxx' and (select count(*) from users) < 10 #")
}
