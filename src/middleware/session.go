package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/cache"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
	"muyu.com/reader_server_go/v1/src/util"
	"net/http"
)

// 不需要登录的路径
var pathSet = map[string]struct{}{
	"/login":    {},
	"/register": {},
}

// SessionCheck 检查 cookie 中的 sessionId, 并且将 session 内容设置到 Context
func SessionCheck(c *gin.Context) {
	path := c.Request.URL.Path
	if _, ok := pathSet[path]; ok {
		c.Next()
		return
	}

	sessionId, err := c.Cookie("sessionId")

	if err != nil || len(sessionId) == 0 {
		logger.Warning("get cookie sessionId error! sessionId = %s, err = %s", sessionId, err)
		returnNeedLogin(c)
		return
	}

	res := cache.Redis().Get(c, "session:"+sessionId)
	byteArr, err := res.Bytes()
	if res.Err() != nil || len(res.Val()) == 0 || err != nil {
		logger.Warning("get session from redis error! sessionId = %s, session = %s, err = %s", sessionId, res.Val(), res.Err())
		returnNeedLogin(c)
		return
	}

	var user ctx.UserCtx
	if err = json.Unmarshal(byteArr, &user); err != nil {
		logger.Warning("unmarshal user json error! session = %s, err = %s", res.Val(), err)
		returnNeedLogin(c)
		return
	}

	ctx.BindUser(c, &user)

	util.GoSafe(func() {
		// 刷新 session 过期时间
		err := service.RefreshSession(c, &user)
		if err != nil {
			logger.Error("[SessionCheck] RefreshSession user = %+v, err = %s", user, err)
		}
	})

}

func returnNeedLogin(c *gin.Context) {
	c.JSON(http.StatusOK, serializer.NeedLogin)
	c.Abort()
}
