package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
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
	if exist, err := user.CheckIfExist(); !exist {
		if err != nil {
			return -1, err
		}
		query := `INSERT IGNORE  user (UserName, Password, Email) VALUES (?,?,?)`
		hashedPassword, err := user.GetHashPassword()
		if err != nil {
			return -1, err
		}
		rs, err := SqlDB.Exec(query, user.Username, hashedPassword, user.Email)
		return HandleSQLResponse(rs, err)
	}
	return -1, errors.New("user already exist")
}
func (user *User) CheckIfExist() (bool, error) {
	query := "SELECT id FROM user WHERE UserName = (?) OR Email = (?) LIMIT 1"
	rows, err := SqlDB.Query(query, user.Username, user.Email)
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
		return true, rows.Close()
	}
	return false, rows.Close()
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

func (user *User) LoginUser() (int64, error) {
	query := ""
	var rows *sql.Rows
	var err error
	if user.Username != "" && user.Email == "" {
		query = "SELECT id,Password FROM user WHERE UserName = (?) LIMIT 1"
		rows, err = SqlDB.Query(query, user.Username)
	} else if user.Email != "" && user.Username == "" {
		query = "SELECT id,Password FROM user WHERE Email = (?) LIMIT 1"
		rows, err = SqlDB.Query(query, user.Email)
	} else {
		return -1, errors.New("username and email must not both be filled")
	}
	if err != nil {
		return -1, err
	}
	var id, hashedPassword string
	for rows.Next() {
		err = rows.Scan(&id, &hashedPassword)
		if err != nil {
			return -1, err
		}
	}
	if hashedPassword == "" || id == "" {
		return -1, errors.New("username/email or password doesn't exist")
	}
	idInt64, _ := strconv.ParseInt(id, 10, 64)
	return idInt64, user.CheckPassword(hashedPassword)
}

func (user *User) StoreToken(token string) (int64, error) {
	query := ""
	if user.Username != "" && user.Email == "" {
		query = "UPDATE user SET user.Token = ? WHERE UserName = ?"
		return HandleSQLResponse(SqlDB.Exec(query, token, user.Username))
	} else if user.Email != "" && user.Username == "" {
		query = "UPDATE user SET user.Token = ? WHERE Email = ?"
		return HandleSQLResponse(SqlDB.Exec(query, token, user.Email))
	} else {
		return -1, errors.New("username and email must not both be filled")
	}
}

func (user *User) Logout() (int64, error) {
	if user.Id == 0 {
		return -1, errors.New("id does not find")
	}
	query := "UPDATE user SET user.Token = '' WHERE id = ?"
	return HandleSQLResponse(SqlDB.Exec(query, user.Id))
}

func (user *User) GetToken() (string, error) {
	query := ""
	result := ""
	if user.Id != 0 {
		query = "SELECT Token FROM user WHERE id = ? LIMIT 1"
		rows, err := SqlDB.Query(query, user.Id)
		if err != nil {
			return "", err
		}
		for rows.Next() {
			err = rows.Scan(&result)
			if err != nil {
				return "", err
			}
		}
		if result != "" {
			return result, rows.Close()
		}
		return "", rows.Close()
	}
	return "", errors.New("id does not match")
}
