package ctx

import (
	"github.com/gin-gonic/gin"
)

type UserCtx struct {
	ID int

	// 账号
	Account string

	// 用户角色
	Role int

	// 客户端类型
	ClientType int

	SessionId string
}

func BindUser(c *gin.Context, user *UserCtx) {
	c.Set("user", user)
}

func GetUser(c *gin.Context) *UserCtx {
	res, _ := c.Get("user")

	return res.(*UserCtx)
}
