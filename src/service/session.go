package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/cache"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/util"
	"time"
)

const (
	SessionTime   = util.Month
	SessionPrefix = "session:"

	SessionListTime      = SessionTime + time.Hour
	sessionListKeyFormat = "sessionList:%d:%d"

	// SessionMaxCount 每个账号每个客户端保留的最大 session 数量
	SessionMaxCount = 3
)

func ParseSessionListKey(userId int, clientType int) string {
	return fmt.Sprintf(sessionListKeyFormat, userId, clientType)
}

// SetSession 设置 session 并且维护活跃的 session 数量, 登出多余的 session
func SetSession(user *model.User, c *gin.Context, clientType int) (sessionId string, err error) {
	sessionId, err = util.RandomID()
	if err != nil {
		logger.Error("Login get session id error, user = %v , err = %s", user, err)
		return
	}

	str, err := json.Marshal(ctx.UserCtx{
		ID:         user.ID,
		Account:    user.Account,
		SessionId:  sessionId,
		Role:       user.Role,
		ClientType: clientType,
	})

	if err != nil {
		logger.Error("[SetSession] json parse error user = %s, err = %s", user, err)
		return
	}

	err = maintainSessionList(c, user, sessionId, clientType)
	if err != nil {
		logger.Error("[SetSession] maintainSessionList user = %+v, err = %s", user, err)
		return
	}

	err = cache.Redis().Set(c, SessionPrefix+sessionId, string(str), SessionTime).Err()
	if err != nil {
		logger.Error("[setSession] redis set error, user = %v , err = %s", user, err)
		return
	}

	return
}

func maintainSessionList(c *gin.Context, user *model.User, sessionId string, clientType int) error {
	sessionListKey := ParseSessionListKey(user.ID, clientType)

	length, err := cache.Redis().RPush(c, sessionListKey, sessionId).Result()
	if err != nil {
		logger.Error("[SetSession] get session list length error, user = %v, err = %s", user, err)
		return err
	}

	err = cache.Redis().Expire(c, sessionListKey, SessionListTime).Err()
	if err != nil {
		logger.Error("[maintainSessionList] redis expire error! sessionListKey = %s, err = %s", sessionListKey, err)
		return err
	}

	if length < SessionMaxCount {
		return nil
	}

	// 设置当前 session 后, session 数量超过限制, 删除最早的 session
	delSessionId, err := cache.Redis().LPop(c, sessionListKey).Result()
	if err != nil {
		logger.Error("[maintainSessionList] LPop error, sessionListKey = %s, err = %s", sessionListKey, err)
		return err
	}
	err = cache.Redis().Del(c, delSessionId).Err()
	if err != nil {
		logger.Error("[maintainSessionList] Del error, delSessionId = %s, err = %s", delSessionId, err)
		return err
	}

	return nil
}

func RefreshSession(c context.Context, userCtx *ctx.UserCtx) error {
	err := cache.Redis().Expire(c, SessionPrefix+userCtx.SessionId, SessionTime).Err()
	if err != nil {
		logger.Error("[RefreshSession] redis expire session error, userCtx = %v, err = %s", userCtx, err)
		return err
	}

	sessionListKey := ParseSessionListKey(userCtx.ID, userCtx.ClientType)
	err = cache.Redis().Expire(c, sessionListKey, SessionListTime).Err()

	if err != nil {
		logger.Error("[RefreshSession] redis expire sessionMap error, userCtx = %s, err = %s", userCtx, err)
		return err
	}
	return nil
}
