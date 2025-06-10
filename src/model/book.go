package model

import (
	"errors"
	"gorm.io/gorm"
)

type Book struct {
	ID int `gorm:"primaryKey"`

	// 书名
	BookName string

	// 作者
	Author string

	// 总页数
	TotalPage int

	// 封面图片地址
	CoverPic string

	// 描述
	Description string

	// 出版社
	Publisher string

	// epubId
	EpubId string

	// 大图封面
	BigCoverPic string
}

type BookAndTag struct {
	Book

	TagId int
}

func (*Book) TableName() string {
	return "book"
}

type BookRepo struct {
	db *gorm.DB
}

func NewBookRepo(db *gorm.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (c *BookRepo) GetBook(id int) (res Book, error error) {
	query := c.db.First(&res, id)
	if query.Error != nil {
		return res, query.Error
	}
	return res, nil
}

// SearchBook 搜索书籍列表
// search 搜索字符串
func (c *BookRepo) SearchBook(search string, tag int, limit int, offset int) ([]BookAndTag, error) {

	query := c.buildSearchQuery(search, tag).Limit(limit).Offset(offset)

	res := make([]BookAndTag, 0, limit)
	dbRes := query.Scan(&res)

	return res, dbRes.Error
}

func (c *BookRepo) buildSearchQuery(search string, tag int) *gorm.DB {

	query := c.db.Table("book").
		Select(" book.id AS id, book_name, total_page, cover_pic, big_cover_pic, COALESCE(tag_id, -1) AS tag_id ").
		Joins(" LEFT JOIN book_tag_rel ON book.id = book_tag_rel.book_id ")

	if tag != -1 {
		query = query.Where(" tag_id = ? ", tag)
	}

	if len(search) != 0 {
		query = query.Where(" ( book_name LIKE CONCAT('%', ?, '%') OR author LIKE CONCAT('%', ?, '%') ) ", search, search)
	}

	return query
}

func (c *BookRepo) CountSearchBook(search string, tag int) int {
	var res int64
	query := c.buildSearchQuery(search, tag)
	query.Count(&res)
	return int(res)
}

// CheckBooId 检查 bookId 是否存在
func (c *BookRepo) CheckBooId(bookId int) (bool, error) {

	var res Book
	query := c.db.First(&res, bookId)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if query.Error != nil {
		return false, query.Error
	}

	return true, nil
}
