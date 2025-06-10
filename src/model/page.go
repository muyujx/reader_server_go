package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"muyu.com/reader_server_go/v1/src/logger"
	"strconv"
)

type HtmlContent struct {
	ID      int `gorm:"primaryKey"`
	Content string
	BookId  int
}

type PageHtml struct {
	ID            int `gorm:"primaryKey"`
	BookId        int
	PageIdx       int
	PageOfChapter int
	Chapter       int
	ContentId     int // 为空说明没有下载
}

type PageRepo struct {
	db *gorm.DB
}

func NewPageRepo(db *gorm.DB) *PageRepo {
	shardingInit(db)
	return &PageRepo{db: db}
}

func shardingInit(db *gorm.DB) {
	err := db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "book_id",
		NumberOfShards:      10,
		PrimaryKeyGenerator: sharding.PKSnowflake,
		ShardingAlgorithm: func(val any) (string, error) {
			if t, ok := val.(string); ok {
				intV, err := strconv.Atoi(t)
				if err != nil {
					return "", err
				}
				return "_" + strconv.Itoa(intV%10), nil
			}

			if t, ok := val.(int); ok {
				return "_" + strconv.Itoa(t%10), nil
			}

			return "", errors.New(fmt.Sprintf("book_id type is %T", val))
		},
	}, "raw_html"))

	if err != nil {
		logger.Panic("register sharding fail! err = %s", err)
	}
}

type PageContentQueryRes struct {
	PageIndex int
	Content   string
	Title     string
}

type PageQueryContent struct {
	PageIdx   int
	ContentId int `gorm:"column:content_id"`
	Title     string
}

func (c *PageRepo) QueryPage(bookId int, startPage int, pageSize int) ([]PageContentQueryRes, error) {
	sql := ` SELECT content_id, page_idx,
				   ( SELECT label
					FROM contents
					WHERE start_page <= page_html.page_idx
					  AND book_id = page_html.book_id
					ORDER BY start_page DESC
					LIMIT 1 ) title
			FROM page_html
			WHERE page_html.book_id = ?
			  AND page_idx >= ?
			ORDER BY page_idx
			LIMIT ? `

	query := c.db.Raw(sql, bookId, startPage, pageSize)

	var pageArr = make([]PageQueryContent, pageSize)
	query = query.Scan(&pageArr)

	if query.Error != nil {
		return nil, query.Error
	}

	contentIdArr := make([]int, len(pageArr))
	for i, v := range pageArr {
		contentIdArr[i] = v.ContentId
	}

	query = c.db.Table("raw_html").Select("id, content").Where("book_id = ?", bookId).Where("id IN (?)", contentIdArr)
	var contentArr []HtmlContent

	query = query.Scan(&contentArr)
	contentMap := make(map[int]string, len(contentArr))
	for _, v := range contentArr {
		contentMap[v.ID] = v.Content
	}

	if query.Error != nil {
		return nil, query.Error
	}

	var res []PageContentQueryRes
	for _, v := range pageArr {
		res = append(res, PageContentQueryRes{
			Title:     v.Title,
			PageIndex: v.PageIdx,
			Content:   contentMap[v.ContentId],
		})
	}

	return res, nil
}
