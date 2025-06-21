package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
	"net/http"
)

func SearchBookList(c *gin.Context) {
	var svc service.BookListService

	if err := c.ShouldBindJSON(&svc); err == nil {
		c.JSON(http.StatusOK, svc.SearchBookList(c))
	} else {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
	}
}

func SearchOnType(c *gin.Context) {
	var svc service.SearchBookOnTypeSvc

	if err := c.ShouldBindJSON(&svc); err == nil {
		c.JSON(http.StatusOK, svc.SearchBookOnType())
	} else {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
	}
}

func GetBookInfo(c *gin.Context) {
	bookId := c.Query("bookId")
	c.JSON(http.StatusOK, service.GetBookInfo(bookId))
}

func GetContents(c *gin.Context) {
	bookId := c.Query("bookId")
	c.JSON(http.StatusOK, service.GetContents(bookId, c))
}
