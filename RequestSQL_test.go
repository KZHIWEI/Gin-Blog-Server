package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var demoUser = User{
Username: "demo",
Password: "password1",
Email:    "demo@demo.com",
}
func initTest(){
	_, err := demoUser.DeleteUser()
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	LoadEnv()
	err := initSQL()
	if err != nil {
		panic(err)
	}
	initTest()
	code := m.Run()
	//shutdown()
	os.Exit(code)
}

func TestCreateUser(t *testing.T){
	_ ,err := demoUser.CreateUser()
	assert.Nil(t,err)
}

func TestDeduplicateUser(t *testing.T) {
	id ,err := demoUser.CreateUser()
	assert.Equal(t,"user already exist",err.Error())
	assert.EqualValues(t,id,-1)
}