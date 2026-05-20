package model

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"muyu.com/reader_server_go/v1/src/config"
	"muyu.com/reader_server_go/v1/src/logger"
	"time"
)

type dbRepo struct {
	DB *gorm.DB

	UserRepo *UserRepo

	PageRepo *PageRepo

	BookRepo *BookRepo

	ContentsRepo *ContentsRepo

	BookTagRepo *BookTagRepo

	FavoriteBookRepo *FavoriteBookRepo

	ReadHistoryRepo *ReadHistoryRepo
}

var DbOp *dbRepo

func InitDb() {
	mysqlConfig := config.Config.Mysql
	curDb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlConfig.User,
		mysqlConfig.Pwd, mysqlConfig.Host, mysqlConfig.DB))
	if err != nil {
		logger.Panic("mysql connect fail! err = %s", err)
	}

	curDb.SetConnMaxLifetime(time.Minute * 10)
	curDb.SetMaxOpenConns(20)
	curDb.SetMaxIdleConns(10)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: curDb,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		logger.Panic("gorm init fail! err = %s", err)
	}

	if gin.Mode() != gin.ReleaseMode {
		gormDB = gormDB.Debug()
	}

	DbOp = &dbRepo{
		DB:               gormDB,
		UserRepo:         NewUseRepo(gormDB),
		PageRepo:         NewPageRepo(gormDB),
		BookRepo:         NewBookRepo(gormDB),
		ContentsRepo:     NewContentsRepo(gormDB),
		BookTagRepo:      NewBookTagRepo(gormDB),
		FavoriteBookRepo: NewFavoriteBookRepo(gormDB),
		ReadHistoryRepo:  NewReadHistoryRepo(gormDB),
	}

}
