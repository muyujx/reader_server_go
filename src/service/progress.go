package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/util"
	"time"
)

type AddReadingTimeSvc struct {
	BookId      int `json:"bookId"`
	ReadingCost int `json:"readingCost"`
	Page        int `json:"page"`
}

// UpdateReadingProgress 更新阅读进度
func (svc *AddReadingTimeSvc) UpdateReadingProgress(c *gin.Context) serializer.Response {
	user := ctx.GetUser(c)
	bookId := svc.BookId
	// 阅读消耗时间
	readingCost := svc.ReadingCost
	// 当前阅读到的页数
	page := svc.Page

	if readingCost <= 0 {
		return serializer.ParamErrMsg("readingCost must greater than 0")
	}
	if page <= 0 {
		return serializer.ParamErrMsg("page must greater than 0")
	}
	fBook, err := model.DbOp.FavoriteBookRepo.GetFavoriteBook(user.ID, bookId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return serializer.ErrMsgRes("this book not favorite book!")
	} else if err != nil {
		logger.Error("[UpdateReadingProgress] GetFavoriteBook err: %v", err)
		return serializer.UnknownErr
	}

	curTime := time.Now().UnixMilli()
	err = model.DbOp.DB.Transaction(func(tx *gorm.DB) (errTx error) {
		errTx = updateReadHistory(user.ID, bookId, readingCost, page, fBook.Page)
		if errTx != nil {
			logger.Error("[UpdateReadingProgress] updateReadHistory err: %v", err)
			return
		}
		errTx = model.DbOp.FavoriteBookRepo.UpdateUserReadingProgress(user.ID,
			bookId, readingCost, page, curTime)
		if errTx != nil {
			logger.Error("[UpdateReadingProgress] AddBookReadingTime err: %v", err)
			return
		}
		return
	})

	if err != nil {
		logger.Error("[UpdateReadingProgress] Update Transaction fail, err: %v", err)
		return serializer.UnknownErr
	}
	return serializer.Success
}

// createReadHistory 创建当天的阅读历史记录
func createReadHistory(userId int, bookId int,
	dayStr string, page int, readCost int,
	curPage int) (*model.ReadHistory, error) {
	curTime := time.Now().Unix()
	record := &model.ReadHistory{
		UserId:      userId,
		BookId:      bookId,
		DayStr:      dayStr,
		StartPage:   curPage,
		EndPage:     util.MaxInt(curPage, page),
		ReadingCost: readCost,
		CreateTime:  curTime,
		UpdateTime:  curTime,
	}
	err := model.DbOp.ReadHistoryRepo.AddReadHistory(record)

	return record, err
}

// updateReadHistory 更新阅读历史记录
// fBookPage 当前阅读到的页数
func updateReadHistory(
	userId int,
	bookId int,
	readingCost int,
	page int,
	fBookPage int) error {

	dayStr := time.Now().Format("2006-01-02")
	historyExists, err := model.DbOp.ReadHistoryRepo.CheckExists(userId, bookId, dayStr)
	if err != nil {
		logger.Error("[UpdateReadingProgress] ReadHistoryRepo.CheckExists err: %v", err)
		return err
	}

	if historyExists {
		// 当天有阅读记录, 更新阅读进度和阅读时间
		err = model.DbOp.ReadHistoryRepo.UpdateReadingPageCost(userId, bookId, dayStr, readingCost, page)
		if err != nil {
			logger.Error("[UpdateReadingProgress] ReadHistoryRepo.UpdateReadingPageCost err: %v", err)
			return err
		}
	} else {
		// 当天还没有阅读记录
		_, err = createReadHistory(userId, bookId, dayStr, page, readingCost, fBookPage)
		if err != nil {
			logger.Error("[UpdateReadingProgress] createReadHistory err: %v", err)
			return err
		}
	}
	return nil
}

type ListReadHistorySvc struct {
	StartStr string `json:"startStr"`
	EndStr   string `json:"endStr"`
}

type ListReadHistoryRes struct {
	DayStr string `json:"dayStr"`

	BookList []*BookReadHistory `json:"bookList"`
}

type BookReadHistory struct {

	// 书名
	BookName string `json:"bookName"`

	// 封面图片 url
	CoverPic string `json:"coverPic"`

	// 阅读时间
	ReadingCost int `json:"readingCost"`

	// 起始页
	StartPage int `json:"startPage"`

	// 结束页
	EndPage int `json:"endPage"`
}

// ListReadHistory 获取历史阅读记录
func (svc *ListReadHistorySvc) ListReadHistory(c *gin.Context) serializer.Response {

	startStr := svc.StartStr
	endStr := svc.EndStr
	startTime, err := time.Parse(util.DateDayFormat, startStr)
	if err != nil {
		return serializer.ParamErrMsg("start date str is invalid!")
	}
	startStr = startTime.Format(util.DateDayFormat)
	endTime, err := time.Parse(util.DateDayFormat, endStr)
	if err != nil {
		return serializer.ParamErrMsg("end date str is invalid!")
	}
	endStr = endTime.Format(util.DateDayFormat)

	user := ctx.GetUser(c)
	list, err := model.DbOp.ReadHistoryRepo.ListReadHistory(user.ID, startStr, endStr)
	if err != nil {
		logger.Error("[ListReadHistory] ReadHistoryRepo.ListReadHistory error! %v", err)
	}

	res := make([]*ListReadHistoryRes, 0, len(list))
	for _, item := range list {
		dayStr := item.DayStr

		coverPic := item.CoverPic
		if len(coverPic) == 0 {
			coverPic = item.BigCoverPic
		}
		bookListItem := &BookReadHistory{
			BookName:    item.BookName,
			CoverPic:    coverPic,
			ReadingCost: item.ReadingCost,
			StartPage:   item.StartPage,
			EndPage:     item.EndPage,
		}

		if len(res) == 0 || res[len(res)-1].DayStr != dayStr {
			curItem := &ListReadHistoryRes{
				DayStr:   dayStr,
				BookList: []*BookReadHistory{bookListItem},
			}
			res = append(res, curItem)
		} else {
			bookList := res[len(res)-1].BookList
			bookList = append(bookList, bookListItem)
			res[len(res)-1].BookList = bookList
		}
	}

	return serializer.SuccessRes(res)
}
