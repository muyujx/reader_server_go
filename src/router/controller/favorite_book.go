package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
	"net/http"
	"strconv"
)

func AddFavoriteBook(c *gin.Context) {
	svc := &service.AddFavoriteBookSvc{}

	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.AddFavoriteBook(c))
	} else {
		serializer.ReturnJson(c, serializer.ParamErr)
	}
}

func DeleteFavoriteBook(c *gin.Context) {
	svc := &service.DeleteFavoriteBookSvc{}

	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.DeleteFavoriteBook(c))
	} else {
		serializer.ReturnJson(c, serializer.ParamErr)
	}
}

func ListFavoriteBook(c *gin.Context) {
	svc := &service.FavoriteListSvc{}
	if err := c.ShouldBindJSON(svc); err == nil {
		serializer.ReturnJson(c, svc.GetFavoriteBookList(c))
	} else {
		serializer.ReturnJson(c, serializer.ParamErr)
	}
}

func UpdateUserCurrentPage(c *gin.Context) {
	var svc service.UpdateUsrCurPageSvc

	if err := c.ShouldBindJSON(&svc); err == nil {
		c.JSON(http.StatusOK, svc.UpdateUsrCurPage(c))
	} else {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
	}
}

func GetUserCurrentPage(c *gin.Context) {
	bookId, err := strconv.Atoi(c.Query("bookId"))

	if err != nil {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
		return
	}

	c.JSON(http.StatusOK, service.GetUserCurPage(c, bookId))
}
