package main

import (
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)
func Authenticator(c *gin.Context) (i interface{}, err error) {
	var loginValues User
	if err := c.ShouldBind(&loginValues); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	return AuthorizeLogin(&loginValues)
	//if err != nil {
	//	return nil,err
	//}
	//return payload , nil
}
func PayloadFunc(data interface{}) jwt.MapClaims {
	return jwt.MapClaims {
		"UserName" : data.(*UserPayLoad).UserName,
		"id" : data.(*UserPayLoad).id,
		"Email" : data.(*UserPayLoad).Email,
	}
}
func AuthMiddleware(key string) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:                 "Blog-server",
		Key:                   []byte(key),
		Authenticator:         Authenticator,
		LoginResponse:         LoginResponse,
		PayloadFunc:           PayloadFunc,
		LogoutResponse:        nil,
		RefreshResponse:       nil,
		IdentityHandler:       nil,
		IdentityKey:           "",
		TokenLookup:           "",
		TokenHeadName:         "",
		TimeFunc:              nil,
		HTTPStatusMessageFunc: nil,
		PrivKeyFile:           "",
		PubKeyFile:            "",
		SendCookie:            false,
		SecureCookie:          false,
		CookieHTTPOnly:        false,
		CookieDomain:          "",
		SendAuthorization:     false,
		DisabledAbort:         false,
		CookieName:            "",
	})
	return authMiddleware, err
}
