package service

import (
	"muyu.com/reader_server_go/v1/src/model"
	"muyu.com/reader_server_go/v1/src/serializer"
)

type PageQueryService struct {
	BookId    int `json:"bookId"`
	StartPage int `json:"startPage"`
	PageSize  int `json:"pageSize"`
}

type PageQueryRes struct {
	Content    string `json:"content"`
	Page       int    `json:"page"`
	Title      string `json:"title"`
	TopChapter int    `json:"topChapter"`
}

func (svc *PageQueryService) QueryPage() serializer.Response {
	bookId := svc.BookId
	startPage := svc.StartPage
	pageSize := svc.PageSize

	// 最多查询 30 页
	if pageSize > 30 {
		pageSize = 30
	}
	var res []PageQueryRes
	if pageSize <= 0 {
		return serializer.SuccessRes(res)
	}

	pageArr, err := model.DbOp.PageRepo.QueryPage(bookId, startPage, pageSize)
	if err != nil {
		return serializer.ErrRes(err)
	}

	for _, v := range pageArr {
		res = append(res, PageQueryRes{
			Content: v.Content,
			Page:    v.PageIndex,
			Title:   v.Title,
		})
	}

	return serializer.SuccessRes(res)
}
