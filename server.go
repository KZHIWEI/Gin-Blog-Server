package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
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
	}
}

type Config struct {
	JwtToken   string
	DbAddress  string
	Port       string
	DbPassword string
}

var GlobalConfig Config

func main() {
	LoadEnv()
	if err := initSQL(); err != nil {
		log.Fatal("Not able to connect SQL server")
		return
	}
	r := gin.Default()
	auth, err := AuthMiddleware(GlobalConfig.JwtToken)
	if err != nil {
		panic(err.Error())
	}
	r.POST("/login", auth.LoginHandler)
	r.POST("/register", RegisterHandler)
	api := r.Group("/api")
	api.Use(auth.MiddlewareFunc())
	{

	}
	log.Fatal(r.Run(GlobalConfig.Port))
	SqlDB.Close()
}
