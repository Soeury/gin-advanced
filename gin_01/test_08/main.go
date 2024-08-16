package main

import (
	"database/sql/driver"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// sqlx 实现批量插入   好喜欢啊啊啊啊啊啊啊啊啊

// 占位符 ? 在内部其实是一个 bindvars (查询占位符)
//  ---- bindvars 仅用来在sql语句中插入值，不允许更改sql语句的结构，如  select ? , ? from ?  这是不可以的

var db *sqlx.DB

// 数据库连接
func initMysql() (err error) {

	dsn := "root:123456@tcp(localhost:3306)/my_2?charset=utf8mb4&parseTime=True"
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("found err in sqlxConnect : %v\n", err)
		return
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(200)
	return
}

type Users struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func (u Users) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

// sqlx.In插入  :  拼接语句和参数, 实现数据批量插入, 注意传入的参数是 []interface{}
func BantchInsertRow(users []interface{}) (err error) {

	// 这里要插入几个数据就写几个(?)
	sqlStr := "insert into users (name , age) values (?) , (?) , (?)"
	query, args, err := sqlx.In(sqlStr, users...)
	if err != nil {
		fmt.Printf("found err in sqlx.In : %v\n", err)
		return err
	}

	fmt.Println(query)
	fmt.Println(args)
	_, err = db.Exec(query, args...)
	return err
}

// sqlx.In查询   :  根据给定的一个 ID 集合查询对应的数据
func BantchQueryRow(ids []int) (users []Users, err error) {

	sqlStr := "select id , name , age from users where id in (?)"
	query, args, err := sqlx.In(sqlStr, ids)
	if err != nil {
		fmt.Printf("found err in sqlx.In : %v\n", err)
		return
	}

	//  * sqlx.In 返回带 ? 的 Bindvar查询语句，用 Rebind重新绑定
	query = db.Rebind(query)
	err = db.Select(&users, query, args...)
	return users, err
}

// sqlx.In 按照指定id的顺序进行数据的返回
func QueryReturnByOrders(ids []int) (users []Users, err error) {

	sqlStr := "select * from users where id in (?) order by find_in_set (id , ?)"
	strIDs := make([]string, 0, len(ids))
	for _, v := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", v))
	}

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(strIDs, ","))
	if err != nil {
		fmt.Printf("found err in sqlx.In : %v\n ", err)
		return
	}

	query = db.Rebind(query)
	err = db.Select(&users, query, args...)
	return users, err
}

func main() {

	err := initMysql()
	if err != nil {
		fmt.Printf("found err in initMysql : %v\n", err)
		return
	}
	defer db.Close()
	fmt.Println("connect with my_2 successed!")

	// sqlx.In 批量插入数据
	u1 := Users{Name: "xx", Age: 18}
	u2 := Users{Name: "xxx", Age: 28}
	u3 := Users{Name: "xxxx", Age: 38}
	s := []interface{}{u1, u2, u3}
	BantchInsertRow(s)

	// sqlx.In 查询
	var users []Users
	users, err = BantchQueryRow([]int{6, 7, 8})
	if err != nil {
		fmt.Printf("found err in bantchQueryRow : %v\n", err)
		return
	}

	for _, v := range users {
		fmt.Printf("%+v\n", v)
	}
	fmt.Println()

	// query return by orders 按照id顺序返回查询结果，一般不会这么做，直接用一个for循环让id按照指定顺序打印即可
	var users2 []Users
	users2, err = QueryReturnByOrders([]int{7, 6, 8, 1})
	if err != nil {
		fmt.Printf("found err in bantchQueryRow : %v\n", err)
		return
	}

	for _, v := range users2 {
		fmt.Printf("%+v\n", v)
	}
}
