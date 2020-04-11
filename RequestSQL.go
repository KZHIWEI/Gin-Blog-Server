package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var SqlDB *sql.DB

func initSQL() error {
	var err error
	SQLAddress := fmt.Sprintf("root:%s@tcp(%s)/gin-test", GlobalConfig.DbPassword, GlobalConfig.DbAddress)
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

func HandleSQLResponse(rs sql.Result, err error) (int64, error) {
	if err != nil {
		return -1, err
	}
	id, err := rs.LastInsertId()
	if err != nil {
		return id, err
	}
	return id, nil
}
func (user *User) CreateUser() (int64, error) {
	if exist , err := user.CheckIfExist();!exist {
		if err != nil {
			return -1,err
		}
		query := `INSERT IGNORE  user (UserName, Password, Email) VALUES (?,?,?)`
		hashedPassword , err := user.GetHashPassword()
		if err != nil {
			return -1 ,err
		}
		rs, err := SqlDB.Exec(query, user.Username, hashedPassword, user.Email)
		return HandleSQLResponse(rs, err)
	}
	return -1 , errors.New("user already exist")
}
func (user *User) CheckIfExist()(bool, error){
	//query := "SELECT * FROM user WHERE user.UserName= ? OR user.Email= ?"
	query := "SELECT id FROM user WHERE UserName = (?) OR Email = (?)"
	rows, err := SqlDB.Query(query,user.Username,user.Email)
	if err != nil {
		return false, err
	}
	result := ""
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return false, err
		}
	}
	if result != "" {
		return true ,rows.Close()
	}
	return false ,rows.Close()
}

func (user *User) DeleteUser() (int64, error) {
	if user.Username != "" {
		rs, err := SqlDB.Exec("DELETE FROM user WHERE user.UserName = (?)", user.Username)
		return HandleSQLResponse(rs, err)
	} else if user.Email != "" {
		rs, err := SqlDB.Exec("DELETE FROM user WHERE user.Email = (?)", user.Email)
		return HandleSQLResponse(rs, err)
	}
	return -1, errors.New("empty username or email")
}

func (user *User) LoginUser()(int64, error){

}
