package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserPayLoad struct {
	id       string
	UserName string
	Email    string
}

type User struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email"`
}

func (user User)String()string{
	return fmt.Sprintf("UserName: %s , Email: %s , Password: %s ",user.Username,user.Email,user.Password)
}

func (user *User) GetHashPassword() (string,error) {
	if len(user.Password) == 0 {
		return "",errors.New("password should not be empty")
	}
	bytePassword := []byte(user.Password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(passwordHash),err
}
func (user *User) checkPassword(hashedPassword string) error {
	bytePassword := []byte(user.Password)
	byteHashedPassword := []byte(hashedPassword)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
