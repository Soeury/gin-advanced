package main

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// sqlx 的基本使用
var db *sqlx.DB

func initDB() (err error) {

	dsn := "root:123456@tcp(localhost:3306)/my_2?charset=utf8mb4&parseTime=True"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("found err in sqlx.connect : %v\n", err)
		return
	}

	db.SetMaxIdleConns(200)
	db.SetMaxIdleConns(10)
	return
}

type Users struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

// 单行查询
func getRow() {
	sqlStr := "select id , name , age from users where id = ?"
	var u Users
	err := db.Get(&u, sqlStr, 1) // 注意这里是 db.Get 还是 sqlx.Get
	if err != nil {
		fmt.Printf("found err in db.Get : %v\n", err)
		return
	}
	fmt.Printf("%+v\n", u)
}

// 多行查询
func getMultiRow() {
	sqlStr := "select id , name , age from users where id > ?"
	var u []Users
	err := db.Select(&u, sqlStr, 0)
	if err != nil {
		fmt.Printf("found err in db.Select : %v\n", err)
		return
	}

	if len(u) == 0 {
		fmt.Printf("not found")
		return
	}

	for _, v := range u {
		fmt.Printf("%+v\n", v)
	}
}

// 插入
func insertRow() {

	sqlStr := "insert into users(name , age) values(?,?)"
	ret, err := db.Exec(sqlStr, "iii", 9)
	if err != nil {
		fmt.Printf("found err in db.Exec : %v\n", err)
		return
	}
	num, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("found err in ret.lastinsertid : %v\n", err)
		return
	}
	fmt.Printf("inset seccessed , the last insert id is : %d", num)
}

// 更新
func updateRow() {

	sqlStr := "update users set age = 1 where id = ?"
	ret, err := db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Printf("found err in db.Exec : %v\n", err)
		return
	}

	num, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("found err in ret.rowsaffected : %v\n", err)
		return
	}
	fmt.Printf("update successed , num of rows affected is : %d", num)
}

// 删除
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
	fmt.Printf("delete successed , num of rows affected id : %d", num)
}

// nameExec : 当sql语句的占位符过多的时候, 可能出现传参错的的现象, 使用nameExec时 可以使用 :name 的形式作占位符
func insert_nameExec() (err error) {
	sqlStr := "insert into users(name , age) values (:name , :age)"
	m1 := map[string]interface{}{
		"name": "qqq",
		"age":  5,
	}

	// _ , err := db.NameExec(sqlStr , map)
	_, err = db.NamedExec(sqlStr, m1)
	if err != nil {
		fmt.Printf("found err in db.NmaeExec : %v\n", err)
		return
	}
	return
}

// namedQuery : 使用 map 命名查询
func query_namedQuery() (err error) {
	sqlStr := "select * from users where name = :name"
	m := map[string]interface{}{
		"name": "iii",
	}
	rows, err := db.NamedQuery(sqlStr, m)
	if err != nil {
		fmt.Printf("found err in db.NameQuery : %v\n", err)
		return
	}
	defer rows.Close()

	// 注意下面使用 structScan 把对象映射到一个结构体里面
	for rows.Next() {
		var u Users
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("found err in rows.scan : %v\n", err)
			continue
		}
		fmt.Printf("%+v\n", u)
	}
	return
}

// 事务操作
func transaction() (err error) {

	tx, err := db.Beginx()
	if err != nil {
		fmt.Printf("found err in db.Begin : %v\n", err)
		return err
	}

	// 这个 defer 是重点
	defer func() {
		p := recover()
		if p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback()
		} else {
			fmt.Println("commit!")
			err = tx.Commit()
		}
	}()

	// 第一个更新
	sqlStr1 := "update users set age = 222 where id = ?"
	ret1, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		return err
	}

	num1, err := ret1.RowsAffected()
	if err != nil {
		return err
	}
	if num1 != 1 {
		return errors.New("rows affected num mistaken")
	}

	// 第二个更新
	sqlStr2 := "update users set age = 333 where id = ?"
	ret2, err := tx.Exec(sqlStr2, 2)
	if err != nil {
		return err
	}

	num2, err := ret2.RowsAffected()
	if err != nil {
		return err
	}
	if num2 != 1 {
		return errors.New("rows affected num mistaken")
	}

	return err
}

func main() {

	err := initDB()
	if err != nil {
		fmt.Printf("found err in initDB : %v\n", err)
		return
	}
	defer db.Close()
	fmt.Println("connect with my_2 successed!")

	// 单行查询
	getRow()
	fmt.Println()

	// 多行查询
	getMultiRow()
	fmt.Println()

	// 插入
	insertRow()
	fmt.Println()

	// 更新
	updateRow()
	fmt.Println()

	//删除
	deleteRow()
	fmt.Println()

	// nameExec 形式的插入
	insert_nameExec()
	fmt.Println()

	// namedQuery形式的查询
	query_namedQuery()
	fmt.Println()

	// 事务
	transaction()
	fmt.Println()
}
