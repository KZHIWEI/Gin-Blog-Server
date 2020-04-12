package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AuthorizeLogin(user *User) (*UserPayLoad, error) {
	if user.Email == "" && user.Username == "" {
		return nil, errors.New("email or username can not be empty")
	}
	if user.Password == "" {
		return nil, errors.New("password can not be empty")
	}
	id64, err := user.LoginUser()
	if err != nil {
		return nil, err
	}
	return &UserPayLoad{
		id:       int(id64),
		UserName: user.Username,
		Email:    user.Email,
	}, nil
}

func RegisterHandler(c *gin.Context) {
	var registerUser User
	err := c.BindJSON(&registerUser)
	if err != nil {
		ResponseError(c, err)
		return
	}
	if registerUser.Email == "" || registerUser.Username == "" {
		err = errors.New("email and username can not be empty")
		ResponseError(c, err)
		return
	}
	if registerUser.Password == "" {
		err = errors.New("password can not be empty")
		ResponseError(c, err)
		return
	}
	err = registerUser.ValidateUser()
	if err != nil {
		ResponseError(c, err)
		return
	}
	id, err := registerUser.CreateUser()
	if err != nil {
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful register your account", "id": id})
}

func ResponseError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func LoginResponse(c *gin.Context, i int, token string, expire time.Time) {
	user, exist := c.Get("user")
	if !exist {
		ResponseError(c, errors.New("user does not store"))
		return
	}
	_, err := user.(*User).StoreToken(token)
	if err != nil {
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    i,
		"expire":  expire.Format(time.RFC3339),
		"token":   token,
		"message": "successful login",
	})
}
