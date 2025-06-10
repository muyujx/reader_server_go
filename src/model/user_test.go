package model

import (
	"fmt"
	"muyu.com/reader_server_go/v1/src/logger"
	"testing"
)

func TestUserGetAll(t *testing.T) {
	users := DbOp.UserRepo.GetAllUser()
	for _, v := range users {
		t.Log(fmt.Sprintf("%+v", v))
	}
}

func TestUserGetOne(t *testing.T) {
	user, err := DbOp.UserRepo.GetUserByAccount("test")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func TestUserNotFound(t *testing.T) {
	account := "aa"
	_, err := DbOp.UserRepo.GetUserByAccount(account)
	if err != nil {
		logger.Error("Login find user error, account = %s , err = %s", account, err)
	}

}
