package service

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/serializer"
)

const tagCacheKey = "book:tag:all"

func GetAllTag() serializer.Response {
	tags, err := getAllTag()
	if err != nil {
		return serializer.ErrRes(err)
	}
	return serializer.SuccessRes(tags)
}

func getAllTag() (res []model.BookTag, err error) {
	tags, err := model.DbOp.BookTagRepo.GetAllTag()
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
