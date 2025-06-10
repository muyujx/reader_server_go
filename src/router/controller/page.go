package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
	"net/http"
)

func GetPage(c *gin.Context) {
	var svc service.PageQueryService

	if err := c.ShouldBindJSON(&svc); err == nil {
		c.JSON(http.StatusOK, svc.QueryPage())
	} else {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
	}
}
