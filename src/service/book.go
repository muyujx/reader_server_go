package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"muyu.com/reader_server_go/v1/src/cache"
	"muyu.com/reader_server_go/v1/src/ctx"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/serializer"
	"muyu.com/reader_server_go/v1/src/util"
)

type BookListService struct {
	Page     int
	PageSize int
	Search   string
	Tag      int
}

type BookListRes struct {
	Total     int           `json:"total"`
	TotalPage int           `json:"totalPage"`
	Content   []BookInfoRes `json:"content"`
}

type BookInfoRes struct {
	BookId      int    `json:"bookId"`
	BookName    string `json:"bookName"`
	TotalPage   int    `json:"totalPage"`
	CoverPic    string `json:"coverPic"`
	BigCoverPic string `json:"bigCoverPic"`
	TagId       int    `json:"tagId"`
	Favorite    bool   `json:"favorite"`
	Description string `json:"description"`
}

func (svc *BookListService) SearchBookList(c *gin.Context) serializer.Response {
	search := svc.Search
	pageSize := svc.PageSize
	page := svc.Page
	tag := svc.Tag
	user := ctx.GetUser(c)

	if pageSize > 30 {
		pageSize = 30
	}

	res := BookListRes{}
	total := model.DbOp.BookRepo.CountSearchBook(search, tag)
	if total == 0 {
		return serializer.SuccessRes(res)
	}

	books, err := model.DbOp.BookRepo.SearchBook(search, tag, pageSize, pageSize*(page-1))
	if err != nil {
		return serializer.ErrRes(err)
	}

	bookIdList := make([]int, 0, len(books))
	for _, book := range books {
		bookIdList = append(bookIdList, book.ID)
	}
	// 检查用户收藏的书籍
	fBookMap, err := model.DbOp.FavoriteBookRepo.CheckFavoriteBook(user.ID, bookIdList)
	if err != nil {
		logger.Error("[SearchBookList] [CheckFavoriteBook] error: %v", err)
		return serializer.UnknownErr
	}

	res.Content = make([]BookInfoRes, len(books), len(books))
	for idx, book := range books {
		_, favorite := fBookMap[book.ID]

		res.Content[idx] = BookInfoRes{
			BookId:      book.ID,
			BookName:    book.BookName,
			TotalPage:   book.TotalPage,
			CoverPic:    book.CoverPic,
			BigCoverPic: book.BigCoverPic,
			TagId:       book.TagId,
			Favorite:    favorite,
			Description: book.Description,
		}
	}

	res.Total = total
	res.TotalPage = (total + pageSize - 1) / pageSize
	return serializer.SuccessRes(res)
}

type GetBookInfoRes struct {
	BookName    string `json:"bookName"`
	BookId      int    `json:"bookId"`
	TotalPage   int    `json:"totalPage"`
	CoverPic    string `json:"coverPic"`
	BigCoverPic string `json:"bigCoverPic"`
}

func GetBookInfo(bookIdStr string) serializer.Response {
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		return serializer.ErrRes(err)
	}

	book, err := getBookInfo(bookId)
	if err != nil {
		return serializer.ErrRes(err)
	}

	bigCp := book.BigCoverPic
	if len(bigCp) == 0 {
		bigCp = book.CoverPic
	}

	return serializer.SuccessRes(GetBookInfoRes{
		BookName:    book.BookName,
		BookId:      book.ID,
		TotalPage:   book.TotalPage,
		BigCoverPic: bigCp,
		CoverPic:    book.CoverPic,
	})
}

var bookInfoCacheKey = "bookInfo:%d"

func getBookInfo(bookId int) (res model.Book, err error) {
	ctx := context.Background()
	key := fmt.Sprintf(bookInfoCacheKey, bookId)

	str, err := cache.Redis().Get(ctx, key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(str), &res)
		return
	}
	res, err = model.DbOp.BookRepo.GetBook(bookId)
	if err != nil {
		return
	}

	bArr, err := json.Marshal(res)
	if err != nil {
		return
	}

	err = cache.Redis().Set(ctx, key, bArr, util.Month).Err()
	return
}

type ContentsRes struct {
	Chapter   int    `json:"chapter"`
	Label     string `json:"label"`
	StartPage int    `json:"startPage"`
	PageCount int    `json:"pageCount"`
	Level     int    `json:"level"`
}

func GetContents(bookIdStr string, c *gin.Context) serializer.Response {
	bookId, err := strconv.Atoi(bookIdStr)
	if err != nil {
		return serializer.ErrRes(err)
	}

	arr, err := getContents(bookId, c)
	if err != nil {
		return serializer.ErrRes(err)
	}

	res := make([]ContentsRes, len(arr), len(arr))
	for i, v := range arr {
		res[i] = ContentsRes{
			Chapter:   v.Chapter,
			Label:     v.Label,
			StartPage: v.StartPage,
			PageCount: v.PageCount,
			Level:     v.Level,
		}
	}

	return serializer.SuccessRes(res)
}

var bookContentsKey = "book:contents:%d"

func getContents(bookId int, c context.Context) (res []model.Contents, err error) {
	key := fmt.Sprintf(bookContentsKey, bookId)
	str, err := cache.Redis().Get(c, key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(str), &res)
		return
	}
	res, err = model.DbOp.ContentsRepo.GetContents(bookId)
	if err != nil {
		return
	}

	bArr, err := json.Marshal(res)
	if err != nil {
		return
	}
	err = cache.Redis().Set(c, key, bArr, util.Month).Err()
	return
}
