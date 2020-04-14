package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	err := c.ShouldBindBodyWith(&registerUser, binding.JSON)
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
		"id":      user.(*User).Id,
		"expire":  expire.Format(time.RFC3339),
		"token":   token,
		"message": "successful login",
	})
}
func GetTokenFromContext(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", jwt.ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", jwt.ErrInvalidAuthHeader
	}
	return parts[1], nil
}
func LogoutHandler(c *gin.Context) {
	userR, exist := c.Get("user")
	if !exist {
		ResponseError(c, errors.New("user does not exist"))
		return
	}
	user := userR.(*User)
	_, err := user.Logout()
	if err != nil {
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successful log out",
	})
}
func Authorizator(data interface{}, c *gin.Context) bool {
	jwtPayloadValue, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return false
	}
	jwtPayload := jwtPayloadValue.(jwt.MapClaims)
	payloadId := int(jwtPayload["id"].(float64))
	if payloadId != 0 {
		token, err := GetTokenFromContext(c)
		if err != nil {
			return false
		}
		var user = User{
			Id: payloadId,
		}
		userToken, err := user.GetToken()
		if err != nil {
			return false
		}
		if userToken == "" {
			if c.Request.URL.Path == "/api/auth/logout" {
				return true
			}
		}
		if userToken == token {
			return true
		}
		return false
	}
	return false
}

func TestTokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "valid token",
	})
}

func MD5(name string) string{
	name = name + time.Now().String()
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}

func MkdirIfNotExist() error {
	if _, err := os.Stat(GlobalConfig.ImageDir); os.IsNotExist(err) {
		err := os.Mkdir(GlobalConfig.ImageDir,0755)
		if err != nil {
			return err
		}
	}
	return nil
}
func StoreImage(file *multipart.FileHeader) (string,error){
	filename := file.Filename
	err := MkdirIfNotExist()
	if err != nil {
		return "",err
	}
	storeName := GlobalConfig.ImageDir + MD5(filename) + filepath.Ext(filename)
	out, err := os.Create(storeName)
	defer out.Close()
	if err != nil {
		return "",err
	}
	incoming,_ :=file.Open()
	_, err = io.Copy(out, incoming)
	if err != nil {
		return "",err
	}
	return MD5(filename) + filepath.Ext(filename),nil
}

func ImageUploadHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		ResponseError(c, err)
		return
	}
	name,err:=StoreImage(file)
	if err != nil {
		ResponseError(c, err)
		return
	}
	c.JSON(200,gin.H{
		"message":"successful upload image",
		"url": GlobalConfig.URL + "/image/" +name,
	})
}

func MultiImageUploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		ResponseError(c, err)
	}
	files := form.File["upload[]"]
	for _, file := range files {
		log.Println(file.Filename)

		// Upload the file to specific dst.
		// c.SaveUploadedFile(file, dst)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
