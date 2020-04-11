package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var SqlDB *sql.DB

func initSQL() error {
	var err error
	SQLAddress := fmt.Sprintf("root:%s@tcp(%s)/gin-test",GlobalConfig.DbPassword,GlobalConfig.DbAddress)
	SqlDB, err = sql.Open("mysql", SQLAddress)
	if err != nil {
		return err
	}
	err = SqlDB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (reg *Register)CreateUser() (int64,error) {
	rs, err := SqlDB.Exec("INSERT into user (UserName, Password,Email) value (?,?,?)", reg.Username, reg.Password,reg.Email)
	if err != nil{
		return -1,err
	}
	id, err := rs.LastInsertId()
	if err != nil{
		return id,err
	}
	return id ,nil
}