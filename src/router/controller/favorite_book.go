package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
	"net/http"
	"strconv"
)

// AddFavoriteBook 添加收藏书籍
func AddFavoriteBook(c *gin.Context) {
	svc := &service.AddFavoriteBookSvc{}

	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.AddFavoriteBook(c))
	} else {
		serializer.ReturnJson(c, serializer.ParamErr)
	}
}

// DeleteFavoriteBook 删除收藏书籍
func DeleteFavoriteBook(c *gin.Context) {
	svc := &service.DeleteFavoriteBookSvc{}

	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.DeleteFavoriteBook(c))
	} else {
		serializer.ReturnJson(c, serializer.ParamErr)
	}
}

// ListFavoriteBook 获取收藏书籍列表
func ListFavoriteBook(c *gin.Context) {
	svc := &service.FavoriteListSvc{}
	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.GetFavoriteBookList(c))
	} else {
		serializer.ReturnJson(c, serializer.ParamErr)
	}
}

// UpdateUserCurrentPage 修改当前阅读进度
func UpdateUserCurrentPage(c *gin.Context) {
	var svc service.UpdateUsrCurPageSvc

	if err := c.ShouldBindJSON(&svc); err == nil {
		c.JSON(http.StatusOK, svc.UpdateUsrCurPage(c))
	} else {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
	}
}

// GetUserCurrentPage 获取当前书籍的阅读进度
func GetUserCurrentPage(c *gin.Context) {
	bookId, err := strconv.Atoi(c.Query("bookId"))

	if err != nil {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
		return
	}

	c.JSON(http.StatusOK, service.GetUserCurPage(c, bookId))
}

// AddReadingCost 增加阅读耗时
func AddReadingCost(c *gin.Context) {
	svc := &service.AddReadingTimeSvc{}
	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.AddReadingTime(c))
		return
	} else {
		logger.Error("param error", err)
		serializer.ParamErrMsg(err.Error())
	}
}
