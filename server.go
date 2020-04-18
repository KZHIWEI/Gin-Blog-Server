package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func LoadEnv() {
	err := godotenv.Load(".env.yaml")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	GlobalConfig = Config{
		JwtToken:   os.Getenv("JWT-TOKEN"),
		DbAddress:  os.Getenv("DB-Address"),
		Port:       os.Getenv("PORT"),
		DbPassword: os.Getenv("DB-PASSWORD"),
		URL:        os.Getenv("BASE-URL") + os.Getenv("PORT"),
		ImageDir:   os.Getenv("IMAGE-DIR"),
	}
}

type Config struct {
	JwtToken   string
	DbAddress  string
	Port       string
	DbPassword string
	URL        string
	ImageDir   string
}

var GlobalConfig Config

func main() {
	LoadEnv()
	if err := initSQL(); err != nil {
		log.Fatal("Not able to connect SQL server")
		return
	}
	defer SqlDB.Close()
	r := gin.Default()
	auth, err := AuthMiddleware(GlobalConfig.JwtToken)
	if err != nil {
		panic(err.Error())
	}
	r.StaticFS("/image/", http.Dir("images"))
	api := r.Group("/api")
	{
		api.POST("/login", auth.LoginHandler)
		api.POST("/register", RegisterHandler)
	}
	authGroup := api.Group("/auth")
	authGroup.Use(auth.MiddlewareFunc())
	{
		authGroup.POST("/logout", LogoutHandler)
		authGroup.POST("/test-token", TestTokenHandler)
		authGroup.POST("/upload-image", ImageUploadHandler)
		authGroup.POST("/posts", NewPostsHandler)
		//authGroup.POST("/upload-multiple-image", MultiImageUploadHandler)
	}
	log.Fatal(r.Run(GlobalConfig.Port))
}
