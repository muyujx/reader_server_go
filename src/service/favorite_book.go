package service

import (
	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/util"
	"time"
)

type AddFavoriteBookSvc struct {
	BookId int `json:"bookId"`
}

// AddFavoriteBook 添加书籍收藏
func (svc *AddFavoriteBookSvc) AddFavoriteBook(c *gin.Context) serializer.Response {
	user := ctx.GetUser(c)
	bookId := svc.BookId

	exists, err := model.DbOp.BookRepo.CheckBooId(bookId)
	if err != nil {
		logger.Error("[AddFavoriteBook] CheckBooId  err:", err)
		return serializer.UnknownErr
	}
	if !exists {
		return serializer.ParamErrMsg("book not exists")
	}

	curTime := time.Now().Unix()
	err = model.DbOp.FavoriteBookRepo.AddFavoriteBook(&model.FavoriteBook{
		BookId:     bookId,
		UserId:     user.ID,
		Page:       1,
		CreateTime: curTime,
		UpdateTime: curTime,
	})
	if err != nil {
		if util.IsDuplicateKeyError(err) {
			return serializer.ParamErrMsg("book already exists")
		}
		logger.Error("[AddFavoriteBook] AddFavoriteBook err:", err)
		return serializer.UnknownErr
	}
	return serializer.SuccessNoRes()
}

type DeleteFavoriteBookSvc struct {
	BookId int `json:"bookId"`
}

// DeleteFavoriteBook 删除收藏的书籍
func (svc *DeleteFavoriteBookSvc) DeleteFavoriteBook(c *gin.Context) serializer.Response {
	user := ctx.GetUser(c)
	bookId := svc.BookId
	err := model.DbOp.FavoriteBookRepo.DelFavoriteBook(user.ID, bookId)
	if err != nil {
		logger.Error("[DeleteFavoriteBook] DeleteFavoriteBook err:", err)
		return serializer.UnknownErr
	}
	return serializer.Success
}

type UpdateUsrCurPageSvc struct {
	BookId int `json:"bookId"`
	Page   int `json:"page"`
}

// UpdateUsrCurPage 更新已收藏书籍的阅读进度
func (svc *UpdateUsrCurPageSvc) UpdateUsrCurPage(c *gin.Context) serializer.Response {
	user := ctx.GetUser(c)
	bookId := svc.BookId
	page := svc.Page

	if page < 1 {
		return serializer.ParamErrMsg("page must greater than 0")
	}
	curTime := time.Now().UnixMilli()
	err := model.DbOp.FavoriteBookRepo.UpdateUserCurPage(user.ID, bookId, page, curTime)
	if err != nil {
		return serializer.ErrRes(err)
	}
	return serializer.SuccessNoRes()
}

// GetUserCurPage 获取收藏书籍的阅读进度
func GetUserCurPage(c *gin.Context, bookId int) serializer.Response {
	user := ctx.GetUser(c)
	res, err := model.DbOp.FavoriteBookRepo.GetFavoriteBook(user.ID, bookId)
	if util.IsNoFoundError(err) {
		return serializer.SuccessRes(-1)
	}
	if err != nil {
		logger.Error("[GetUserCurPage] GetFavoriteBook user = %+v, bookId = %d,  err = %+v", user, bookId, err)
		return serializer.UnknownErr
	}
	return serializer.SuccessRes(res.Page)
}

type FavoriteListSvc struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type FavoriteListRes struct {
	Total     int                      `json:"total"`
	TotalPage int                      `json:"totalPage"`
	Content   []model.BookFavoriteInfo `json:"content"`
}

func (svc *FavoriteListSvc) GetFavoriteBookList(c *gin.Context) serializer.Response {
	user := ctx.GetUser(c)
	userId := user.ID
	pageSize := svc.PageSize
	page := svc.Page

	if pageSize > 30 {
		pageSize = 30
	}
	if pageSize < 1 {
		pageSize = 1
	}

	total, err := model.DbOp.FavoriteBookRepo.CountBook(userId)
	if err != nil {
		logger.Error("[GetFavoriteBookList] CountBook err:", err)
		return serializer.UnknownErr
	}

	res := &FavoriteListRes{
		Total:     total,
		TotalPage: (total + pageSize - 1) / pageSize,
		Content:   make([]model.BookFavoriteInfo, 0),
	}
	if page > res.TotalPage {
		return serializer.SuccessRes(res)
	}

	content, err := model.DbOp.FavoriteBookRepo.GetFavoriteBookList(userId, page, pageSize)
	if err != nil {
		logger.Error("[GetFavoriteBookList] GetFavoriteBookList err:", err)
		return serializer.UnknownErr
	}
	res.Content = content
	return serializer.SuccessRes(res)
}

type AddReadingTimeSvc struct {
	BookId      int `json:"bookId"`
	ReadingCost int `json:"readingCost"`
}

// AddReadingTime 增加书籍阅读时间
func (svc *AddReadingTimeSvc) AddReadingTime(c *gin.Context) serializer.Response {
	user := ctx.GetUser(c)
	bookId := svc.BookId
	readingCost := svc.ReadingCost

	if readingCost <= 0 {
		return serializer.ParamErrMsg("readingCost must greater than 0")
	}
	err := model.DbOp.FavoriteBookRepo.AddBookReadingTime(user.ID, bookId, readingCost)
	if err != nil {
		logger.Error("[AddReadingTime] AddBookReadingTime err:", err)
		return serializer.UnknownErr
	}
	return serializer.Success
}
