package main

import (
	"fmt"
	core "muyu.com/reader_server_go/v1/src"
	"muyu.com/reader_server_go/v1/src/config"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/router"
)

func main() {

	// 初始化
	core.Init()

	r := router.InitRouter()

	port := config.Config.Port
	if port == 0 {
		port = 8080
	}

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Panic("http server error!")
	}
}
