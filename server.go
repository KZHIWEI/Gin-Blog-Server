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
		JwtToken:  os.Getenv("JWT-TOKEN"),
		DbAddress: os.Getenv("DB-Address"),
		Port:      os.Getenv("PORT"),
	}
}

type Config struct {
	JwtToken  string
	DbAddress string
	Port      string
}

var GlobalConfig Config

func main() {
	LoadEnv()
	r := gin.Default()
	auth, err := AuthMiddleware(GlobalConfig.JwtToken)
	if err != nil {
		panic(err.Error())
	}
	r.POST("/login", auth.LoginHandler)
	api := r.Group("/api")
	api.Use(auth.MiddlewareFunc())
	{

	}
	err = r.Run(GlobalConfig.Port)
	if err != nil {
		panic(err.Error())
	}
}
