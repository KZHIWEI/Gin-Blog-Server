package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func MD5(name string) string {
	name = name + time.Now().String()
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}

func MkdirIfNotExist() error {
	if _, err := os.Stat(GlobalConfig.ImageDir); os.IsNotExist(err) {
		err := os.Mkdir(GlobalConfig.ImageDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
func StoreImage(file *multipart.FileHeader) (string, error) {
	filename := file.Filename
	err := MkdirIfNotExist()
	if err != nil {
		return "", err
	}
	if !ValidImageFormat( filepath.Ext(filename)) {
		return "", errors.New("file is not a image type")
	}
	md5FileName := MD5(filename)
	storeName := GlobalConfig.ImageDir + md5FileName + filepath.Ext(filename)
	out, err := os.Create(storeName)
	defer out.Close()
	if err != nil {
		return "", err
	}
	incoming, _ := file.Open()
	_, err = io.Copy(out, incoming)
	if err != nil {
		return "", err
	}
	return md5FileName + filepath.Ext(filename), nil
}

func ImageUploadHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		ResponseError(c, err)
		return
	}
	name, err := StoreImage(file)
	if err != nil {
		ResponseError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"message": "successful upload image",
		"url":     GlobalConfig.URL + "/image/" + name,
	})
}

func ValidImageFormat(ext string)bool {
	switch strings.ToLower(ext) {
	case ".jpg":
	case ".jpeg":
	case ".gif":
	case ".png":
		return true
	default:
		return false
	}
	return false
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
