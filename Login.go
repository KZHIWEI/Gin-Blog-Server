package main

import (
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
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
	if id64 != -1 && id64 != 0 {
		user.Id = int(id64)
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
		"id" : user.(*User).Id,
		"expire":  expire.Format(time.RFC3339),
		"token":   token,
		"message": "successful login",
	})
}

func LogoutResponse(c *gin.Context,i int){
	claimsV,exist := c.Get("JWT_PAYLOAD")
	claimsId := int(claimsV.(jwt.MapClaims)["id"].(float64))
	if !exist {
		ResponseError(c,errors.New("no jwt token"))
		return
	}
	var user UserAuth
	if err := c.ShouldBind(&user); err != nil {
		ResponseError(c,err)
		return
	}
	if user.Id == claimsId {
		_,err:= user.Logout()
		if err != nil {
			ResponseError(c,err)
			return
		}
		c.JSON(i,gin.H{
			"message":"successful logout",
		})
		return
	}
	c.JSON(http.StatusUnauthorized,gin.H{
		"message":"id does not match",
	})
}

