package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
)

func Register(c *gin.Context) {
	var svc service.RegisterSvc
	if err := c.ShouldBindJSON(&svc); err != nil {
		serializer.ReturnJson(c, serializer.ParamErr)
		return
	}
	serializer.ReturnJson(c, svc.Register())
}
