package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserPayLoad struct {
	id       int
	UserName string
	Email    string
}

var validate *validator.Validate

type User struct {
	Username string `form:"username" json:"username" `
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email"`
}

func (user User) String() string {
	return fmt.Sprintf("UserName: %s , Email: %s , Password: %s ", user.Username, user.Email, user.Password)
}

func (user *User) GetHashPassword() (string, error) {
	if len(user.Password) == 0 {
		return "", errors.New("password should not be empty")
	}
	bytePassword := []byte(user.Password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(passwordHash), err
}
func (user *User) CheckPassword(hashedPassword string) error {
	bytePassword := []byte(user.Password)
	byteHashedPassword := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	if err != nil {
		return errors.New("incorrect password")
	}
	return nil
}

func (user *User) ValidateUser() error {
	validate = validator.New()
	if user.Email != "" {
		err := validate.Var(user.Email, "email")
		if err != nil {
			return errors.New("email address incorrect format")
		}
	}
	if user.Username != "" {
		if len(user.Username) <= 4 {
			return errors.New("username too short")
		}
	}
	if len(user.Password) <= 6 {
		return errors.New("password too short")
	}
	return nil
}

//func (user *User) FixNullPointer(){
//	if user.Email == nil {
//
//	}
//}
