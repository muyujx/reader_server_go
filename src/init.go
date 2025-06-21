package src

import (
	"github.com/gin-gonic/gin"
	"io"
	"muyu.com/reader_server_go/v1/src/cache"
	"muyu.com/reader_server_go/v1/src/config"
	"muyu.com/reader_server_go/v1/src/es"
	"muyu.com/reader_server_go/v1/src/model"
	"os"
)

func Init() {

	// 加载配置文件
	config.InitConfig()

	switch config.Config.Mode {
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	if gin.Mode() == gin.ReleaseMode {
		// 生产环境, 日志默认写到 reader.log 文件中
		if gin.Mode() == gin.ReleaseMode {
			f, _ := os.OpenFile("reader.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 755)
			os.Stdout = f
			os.Stderr = f
			gin.DefaultWriter = io.MultiWriter(f)
		}
	}

	// 初始化 gorm
	model.InitDb()

	// 初始化 redis
	cache.InitRedis()

	// 初始化 es
	es.InitEs()

}
