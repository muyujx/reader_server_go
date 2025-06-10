package router

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/middleware"
	"muyu.com/reader_server_go/v1/src/router/controller"
)

func initGin() {
	gin.ForceConsoleColor()
}

func InitRouter() *gin.Engine {
	initGin()
	r := gin.Default()

	// Recovery 中间件会 recover 任何 panic。如果有 panic 的话，会写入 500。
	r.Use(gin.Recovery())
	r.Use(middleware.SessionCheck)

	// login
	r.POST("/login", controller.Login)
	r.POST("/logout", controller.Logout)
	r.POST("/register", controller.Register)
	r.POST("/check_login", controller.CheckLogin)

	// book
	bookGp := r.Group("/book/info")
	bookGp.POST("/search/book", controller.SearchBookList)
	bookGp.GET("/get/book/info", controller.GetBookInfo)
	bookGp.GET("/get/contents", controller.GetContents)

	// page
	pageGp := r.Group("/book/page")
	pageGp.POST("/html/page", controller.GetPage)

	// user cur page
	userPageGp := r.Group("/user/page")
	userPageGp.POST("/update", controller.UpdateUserCurrentPage)
	userPageGp.GET("/get", controller.GetUserCurrentPage)

	// type
	typeGp := r.Group("/book/tag")
	typeGp.GET("/get/all", controller.GetAllTag)
	typeGp.POST("/change", controller.ChangeTag)

	// favorite_book
	favoriteGp := r.Group("/favorite")
	favoriteGp.POST("/add", controller.AddFavoriteBook)
	favoriteGp.POST("/delete", controller.DeleteFavoriteBook)
	favoriteGp.POST("/list", controller.ListFavoriteBook)
	favoriteGp.POST("/read_cost", controller.AddReadingCost)

	// test
	testGp := r.Group("/test")
	testGp.GET("/user", controller.GetCurUser)

	return r
}
