package main

type UserPayLoad struct {
	id       string
	UserName string
	Email    string
}

type Login struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email"`
}