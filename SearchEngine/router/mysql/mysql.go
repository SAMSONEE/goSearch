package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)


func InitDB() *sql.DB{

	dsn := "root:123456@tcp(127.0.0.1:3306)/yourdb"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Printf("db err:%v\n", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("连接数据库失败", err)
	}

	db.SetMaxOpenConns(5)

	return db
}




