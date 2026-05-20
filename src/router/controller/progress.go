package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
)

// UpdateReadingProgress 更新阅读进度
func UpdateReadingProgress(c *gin.Context) {
	svc := &service.AddReadingTimeSvc{}
	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.UpdateReadingProgress(c))
		return
	} else {
		logger.Error("[UpdateReadingProgress] param error %v", err)
		serializer.ReturnJson(c, serializer.ParamErr)
	}
}

// ListReadHistory 获取历史阅读记录
func ListReadHistory(c *gin.Context) {
	svc := &service.ListReadHistorySvc{}
	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.ListReadHistory(c))
	} else {
		logger.Error("ListReadHistory param error %v", err)
		serializer.ReturnJson(c, serializer.ParamErr)
	}
}
