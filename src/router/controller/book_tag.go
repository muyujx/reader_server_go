package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
	"net/http"
)

// GetAllTag 获取所有书籍标签
func GetAllTag(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetAllTag())
}

// ChangeTag 修改书籍标签
func ChangeTag(c *gin.Context) {
	var svc service.UpdateBookTagSvc
	if err := c.ShouldBindJSON(&svc); err == nil {
		c.JSON(http.StatusOK, svc.ChangeTag(c))
	} else {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
	}
}
