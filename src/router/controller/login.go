package controller

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/service"
	"net/http"
)

func Login(c *gin.Context) {
	var svc service.LoginService

	if err := c.ShouldBindJSON(&svc); err == nil {
		c.JSON(http.StatusOK, svc.Login(c))
	} else {
		c.JSON(http.StatusOK, serializer.ErrRes(err))
	}
}

func CheckLogin(c *gin.Context) {
	serializer.ReturnJson(c, service.CheckLogin(c))
}

func Logout(c *gin.Context) {
	serializer.ReturnJson(c, service.Logout(c))
}
