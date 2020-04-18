package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type UserPayLoad struct {
	id       int
	UserName string
	Email    string
}

var validate *validator.Validate

type User struct {
	Id       int    `form:"id" json:"id"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Email    string `form:"email" json:"email"`
}

func (user User) String() string {
	return fmt.Sprintf("UserName: %s , Email: %s , Password: %s ", user.Username, user.Email, user.Password)
}

func (user *User) GetHashPassword() (string, error) {
	if len(user.Password) == 0 {
		return "", errors.New("password should not be empty")
	}
	bytePassword := []byte(user.Password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(passwordHash), err
}
func (user *User) CheckPassword(hashedPassword string) error {
	bytePassword := []byte(user.Password)
	byteHashedPassword := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	if err != nil {
		return errors.New("incorrect password")
	}
	return nil
}

func (user *User) ValidateUser() error {
	validate = validator.New()
	if user.Email != "" {
		err := validate.Var(user.Email, "email")
		if err != nil {
			return errors.New("email address incorrect format")
		}
	}
	if user.Username != "" {
		if len(user.Username) <= 4 {
			return errors.New("username too short")
		}
	}
	if len(user.Password) <= 6 {
		return errors.New("password too short")
	}
	return nil
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
		defer rows.Close()
	} else if user.Email != "" && user.Username == "" {
		query = "SELECT id,Password FROM user WHERE Email = (?) LIMIT 1"
		rows, err = SqlDB.Query(query, user.Email)
		defer rows.Close()
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
