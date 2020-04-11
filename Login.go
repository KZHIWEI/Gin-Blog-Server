package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AuthorizeLogin(user *User) (*UserPayLoad , error){
	if user.Email == "" && user.Username == "" {
		return nil , errors.New("email or username can not be empty")
	}
	if user.Password == "" {
		return nil , errors.New("password can not be empty")
	}
	id64 , err := user.LoginUser()
	if err != nil {
		return nil, err
	}
	return &UserPayLoad{
		id:       int(id64),
		UserName: user.Username,
		Email:    user.Email,
	},nil
}

func RegisterHandler(c *gin.Context){
	var registerUser User
	err := c.BindJSON(&registerUser)
	if registerUser.Email == "" || registerUser.Username == "" {
		err =  errors.New("email and username can not be empty")
	}
	if registerUser.Password == "" {
		err =  errors.New("password can not be empty")
	}
	if err != nil {
		fmt.Fprintln(c.Writer,err)
		return
	}
	id,err := registerUser.CreateUser()
	if err != nil {
		fmt.Fprintln(c.Writer,err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful register your account","id":id})
}

func LoginResponse(c *gin.Context,i int,token string,expire time.Time){
	fmt.Printf("i: %v token: %s\n", i, token)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"expire": expire.Format(time.RFC3339),
		"token":token,
		"message":"successful login",
	})
}