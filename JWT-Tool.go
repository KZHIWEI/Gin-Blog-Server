package main

import (
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	jwtTool "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

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

func extractClaims(tokenStr string) (jwt.MapClaims, error) {
	hmacSecretString := GlobalConfig.JwtToken// Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwtTool.Parse(tokenStr, func(token *jwtTool.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwtTool.MapClaims); ok && token.Valid {
		return jwt.MapClaims(claims), nil
	} else {
		return nil, errors.New("error: claims does not exist")
	}
}

func GetClaims(c *gin.Context) (jwt.MapClaims,error) {
	token,err := GetTokenFromContext(c)
	if err != nil {
		return nil, err
	}
	claims,err:=extractClaims(token)
	if err != nil{
		return nil, err
	}
	return claims,nil
}

