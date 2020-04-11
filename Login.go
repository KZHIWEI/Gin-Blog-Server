package main

import "errors"

func AuthorizeLogin(loginValues User) (interface{} , error){
	if loginValues.Email == "" && loginValues.Username == "" {
		return nil , errors.New("email or username can not be empty")
	}
	if loginValues.Password == "" {
		return nil , errors.New("password can not be empty")
	}
	return nil,nil

}