package service

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/serializer"
	"strings"
	"time"
)

type RegisterSvc struct {
	Account string `json:"account"`

	Password string `json:"password"`
}

const (
	AccountMinLen  = 6
	AccountMaxLen  = 30
	PasswordMinLen = 8
	PasswordMaxLen = 30
)

func (s RegisterSvc) Register() serializer.Response {
	account := strings.TrimSpace(s.Account)
	password := strings.TrimSpace(s.Password)

	if len(account) < AccountMinLen || len(account) > AccountMaxLen {
		return serializer.ErrMsgRes("account length must between 6 and 30")
	}
	if len(password) < PasswordMinLen || len(password) > PasswordMaxLen {
		return serializer.ErrMsgRes("password length must between 8 and 30")
	}

	curMillis := time.Now().UnixMilli()

	user := &model.User{
		Account:    account,
		Password:   EncPassword(password),
		CreateTime: curMillis,
		UpdateTime: curMillis,
		LastLogin:  curMillis,
	}

	if err := model.DbOp.UserRepo.CreateUser(user); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return serializer.AccountExits
		}
		logger.Warning("create user error! user = %v, err = %s", user, err.Error())
		return serializer.UnknownErr
	}

	return serializer.Success
}
