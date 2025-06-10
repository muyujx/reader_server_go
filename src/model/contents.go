package model

import "gorm.io/gorm"

type Contents struct {
	ID        int `gorm:"primaryKey"`
	BookId    string
	Chapter   int
	Label     string
	PageCount int
	StartPage int
	Level     int
}

func (*Contents) TableName() string {
	return "contents"
}

type ContentsRepo struct {
	db *gorm.DB
}

func NewContentsRepo(db *gorm.DB) *ContentsRepo {
	return &ContentsRepo{db: db}
}

func (c *ContentsRepo) GetContents(bookId int) ([]Contents, error) {
	var res []Contents
	query := c.db.Where("book_id = ?", bookId).Find(&res)

	if query.Error != nil {
		return res, query.Error
	}

	return res, nil
}
