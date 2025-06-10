package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"muyu.com/reader_server_go/v1/src/cache"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/util"
	"strings"
	"time"
)

const (
	Web = iota
	Desktop
)

type LoginService struct {
	Account string `json:"account"`

	Password string `json:"password"`

	ClientType int `json:"clientType"`
}

type LoginRes struct {
	Account string `json:"account"`

	UserId int `json:"userId"`

	Role int `json:"role"`
}

func (svc *LoginService) Login(c *gin.Context) serializer.Response {
	account := strings.TrimSpace(svc.Account)
	password := strings.TrimSpace(svc.Password)

	if len(account) == 0 || len(password) == 0 {
		return serializer.ErrMsgRes("account or password is empty!")
	}
	if svc.ClientType != Web && svc.ClientType != Desktop {
		return serializer.ErrMsgRes("client type error!")
	}

	user, err := model.DbOp.UserRepo.GetUserByAccount(account)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return serializer.ErrMsgRes("account not found!")
	}
	if err != nil {
		logger.Error("Login find user error, account = %s , err = %s", account, err)
	}

	// 校验密码
	curPwd := EncPassword(password)
	if curPwd != user.Password {
		return serializer.ErrMsgRes("account or password error!")
	}

	lastLogin := time.Now().UnixMilli()
	user.LastLogin = lastLogin
	// 设置 session
	sessionId, err := SetSession(&user, c, svc.ClientType)
	if err != nil {
		return serializer.UnknownErr
	}

	// 设置 cookie
	c.SetCookie("sessionId", sessionId, int(3*util.Month/time.Second), "/", "", false, true)

	// 修改用户上次登录时间
	util.GoSafe(func() {
		updateLastLoginTime(&user, lastLogin)
	})

	return serializer.SuccessRes(LoginRes{
		UserId:  user.ID,
		Role:    user.Role,
		Account: account,
	})
}

// 更新用户上次登录时间
func updateLastLoginTime(user *model.User, lastLogin int64) {
	err := model.DbOp.UserRepo.UpdateLoginTime(user.ID, lastLogin)
	if err != nil {
		logger.Error("user last login time update error,user = %v, err = %v", user, err)
	}
}

// EncPassword 加密密码
func EncPassword(password string) string {
	pwdByte := sha256.Sum256([]byte(password))
	return hex.EncodeToString(pwdByte[:])
}

// Logout 登出当前 session
func Logout(c *gin.Context) serializer.Response {
	userCtx := ctx.GetUser(c)
	sessionId := userCtx.SessionId
	err := cache.Redis().Del(c, "session:"+sessionId).Err()
	if err != nil {
		logger.Error("logout redis delete session error, user = %v , err = %s", userCtx, err)
		return serializer.UnknownErr
	}
	return serializer.SuccessNoRes()
}

// CheckLogin 检查客户端 session 还是否有效
// 如果有效, 返回和 login 相同的返回结果, 客户端则不再需要调用 login
func CheckLogin(c *gin.Context) serializer.Response {
	userCtx := ctx.GetUser(c)
	err := RefreshSession(c, userCtx)
	if err != nil {
		logger.Error("Check login session error, user = %v , err = %s", userCtx, err)
		return serializer.UnknownErr
	}

	return serializer.SuccessRes(&LoginRes{
		UserId:  userCtx.ID,
		Role:    userCtx.Role,
		Account: userCtx.Account,
	})
}
