package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initTest(){

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


func TestDuplicateUser(t *testing.T) {
	reg := Register{
		Username: "demo",
		Password: "password",
		Email:    "demo@demo.com",
	}
	id ,err := reg.CreateUser()
	assert.Error(t,err)
	assert.Equal(t,id,-1)
}