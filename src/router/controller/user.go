package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/serializer"
)

func GetCurUser(c *gin.Context) {
	user := ctx.GetUser(c)
	serializer.ReturnJson(c, serializer.SuccessRes(user))
}
