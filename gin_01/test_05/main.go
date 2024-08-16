package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// mysql 事务处理 : 在一次事务里面, 要么全部处理成功, 要么出现错误回滚到最开始的状态
//   ---- 开启事务 begin()
//   ---- 必须手动提交 commit()
//   ---- 事务处理遇到问题需要回滚 rollback()

var db *sql.DB

// 连接数据库
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

// mysql事务处理
func transaction() {

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		fmt.Printf("found err in db.begin : %v\n", err)
		return
	}

	// 第一个更新
	sqlStr1 := "update users set age = 66 where id = ?"
	ret1, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		tx.Rollback()
		fmt.Printf("found err in tx.exec : %v\n", err)
		return
	}

	affRows1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Printf("found err in ret1.rowsaffected : %v\n", err)
		return
	}

	// 第二个更新
	sqlStr2 := "update users set age = 66 where id = ?"
	ret2, err := tx.Exec(sqlStr2, 100)
	if err != nil {
		tx.Rollback()
		fmt.Printf("found err in tx.exec : %v\n", err)
		return
	}

	affRows2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Printf("found err in ret2.rowsaffected : %v\n", err)
		return
	}

	// 最后判断
	fmt.Println(affRows1, affRows2)
	if affRows1 == 1 && affRows2 == 1 {
		tx.Commit()
		fmt.Println("commit successed!")
	} else {
		tx.Rollback()
		fmt.Println("rollback successed!")
	}

	fmt.Println("tranaction end...")
}

func main() {

	err := initMysql()
	if err != nil {
		fmt.Printf("found err in initmysql : %v\n", err)
		return
	}
	defer db.Close()
	fmt.Println("connect with my_2 successed!")

	// 处理事务
	transaction()
}
