package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"muyu.com/reader_server_go/v1/src/cache"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/util"
)

const tagCacheKey = "book:tag:all"

func GetAllTag(c *gin.Context) serializer.Response {
	tags, err := getAllTag(c)
	if err != nil {
		return serializer.ErrRes(err)
	}
	return serializer.SuccessRes(tags)
}

func getAllTag(c *gin.Context) (res []model.BookTag, err error) {
	str, err := cache.Redis().Get(c, tagCacheKey).Result()
	if len(str) != 0 {
		err = json.Unmarshal([]byte(str), &res)
		return
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		return
	}

	tags, err := model.DbOp.BookTagRepo.GetAllTag()
	if err != nil {
		return
	}
	cTags, err := json.Marshal(tags)
	if err != nil {
		return
	}
	err = cache.Redis().Set(c, tagCacheKey, cTags, util.Month).Err()
	if err != nil {
		return
	}

	return tags, nil
}

type UpdateBookTagSvc struct {
	BookId int `json:"bookId"`

	TagId int `json:"tagId"`
}

func (svc *UpdateBookTagSvc) ChangeTag(c *gin.Context) serializer.Response {
	userCtx := ctx.GetUser(c)
	if userCtx.Role != model.Role_Admin {
		return serializer.NoPermission
	}

	err := model.DbOp.BookTagRepo.UpdateBookTag(svc.BookId, svc.TagId)
	if err != nil {
		return serializer.ErrRes(err)
	}
	return serializer.SuccessNoRes()
}
