package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql驱动执行 init()
)

// ---- databases/sql 实现数据库连接
// ---- 源码解析
var db *sql.DB

func initMysql() (err error) {

	// dsn : data source name
	// sql.open还未连接数据库, 只是检查连接步骤是否有填写错误
	dsn := "root:123456@tcp(localhost:3306)/my_2"
	// 这里要使用前面定义的全局变量
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// 尝试建立连接
	err = db.Ping()
	if err != nil {
		fmt.Printf("found err in db.Open : %v\n", err)
		return
	}

	// 数值需要根据具体业务来确定
	db.SetConnMaxLifetime(time.Second * 10) // 设置最大存活时间
	db.SetMaxIdleConns(10)                  // 设置最大空闲数
	db.SetMaxOpenConns(200)                 // 设置最大连接数
	return
}

type Users struct {
	Id   int
	Name string
	Age  int
}

// 单行查询
func queryRow() {

	// 注意SQL语句不要写错了
	sqlStr := "select id , name , age from users where Id = ?"
	var u Users

	// 这里 row.Scan 必须在 db.QueryRow之后 , Scan 必须使用, 否则有可能导致堵塞, 数据库连接不会被释放
	row := db.QueryRow(sqlStr, 1)
	err := row.Scan(&u.Id, &u.Name, &u.Age)
	if err != nil {
		fmt.Printf("found err in row.scan : %v\n", err)
		return
	}
	fmt.Printf("ID:%d  Name:%s  Age:%d\n", u.Id, u.Name, u.Age)
}

// 多行查询
func queryMultiRow() {

	sqlStr := "select id , name , age from users where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("found err in db.query : %v\n", err)
		return
	}

	// 非常重要 : 关闭rows释放持有的数据库连接
	defer rows.Close()

	// 循环读取rows中的数据
	for rows.Next() {
		var u Users
		err := rows.Scan(&u.Id, &u.Name, &u.Age)
		if err != nil {
			fmt.Printf("found err in rows.scan : %v\n", err)
			return
		}
		fmt.Printf("id:%d  name:%s  age:%d\n", u.Id, u.Name, u.Age)
	}
}

// 插入数据
func insertRow() {

	sqlStr := "insert into users(name , age) values (? , ?)"
	ret, err := db.Exec(sqlStr, "wang", 27)
	if err != nil {
		fmt.Printf("found err in db.Exec : %v\n", err)
		return
	}

	insert_id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("found err in ret.lastinsertid : %v\n", err)
		return
	}
	fmt.Printf("the last insert id is : %d\n", insert_id)
}

// 更新数据
func updateRow() {

	sqlStr := "update users set age = ? where id = ?"
	ret, err := db.Exec(sqlStr, 66, 2)
	if err != nil {
		fmt.Printf("found err in db.exec : %v\n", err)
		return
	}

	// 操作影响的行数
	num, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("found err in ret.rowsaffected : %v\n", err)
		return
	}
	fmt.Printf("num of affected rows is : %d\n", num)
}

// 删除数据
func deleteRow() {

	sqlStr := "delete from users where id = ?"
	ret, err := db.Exec(sqlStr, 4)
	if err != nil {
		fmt.Printf("found err in db.Exec : %v\n", err)
		return
	}

	num, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("found err in ret.rowsaffected : %v\n", err)
		return
	}
	fmt.Printf("num of affected rows is : %d\n", num)
}

func main() {

	err := initMysql()
	if err != nil {
		fmt.Printf("found err in initMysql : %v\n", err)
		return
	}

	// panic 检查确保 db 不为 nil
	// close() 要写在判断 panic 后面 , 用来释放数据库连接相关的资源
	defer db.Close()
	fmt.Println("connect with my_1 successed!")

	// sql查询
	queryRow()
	fmt.Println("queryRow end...")
	queryMultiRow()
	fmt.Println("queryMultiRow end...")

	// sql插入
	insertRow()
	fmt.Println("insertRow end...")

	// sql更新
	updateRow()
	fmt.Println("updateRow end...")

	// sql删除
	deleteRow()
	fmt.Println("deleteRow end...")
}
