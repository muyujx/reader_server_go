package model

import (
	"errors"
	"gorm.io/gorm"
)

type BookTag struct {
	ID int `gorm:"primaryKey" json:"id"`
	// 类型名
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type BookTagRel struct {
	ID int `gorm:"primaryKey" json:"id"`

	TagID int `json:"tag_id"`

	BookID int `json:"book_id"`
}

type BookTagRepo struct {
	db *gorm.DB
}

func NewBookTagRepo(db *gorm.DB) *BookTagRepo {
	return &BookTagRepo{db: db}
}

func (*BookTag) TableName() string {
	return "book_tag"
}
func (*BookTagRel) TableName() string {
	return "book_tag_rel"
}

func (c *BookTagRepo) GetAllTag() ([]BookTag, error) {
	var res []BookTag
	query := c.db.Model(&BookTag{}).Order("`order`").Find(&res)

	if query.Error != nil {
		return nil, query.Error
	}
	return res, nil
}

// 检查标签是否存在
func (c *BookTagRepo) checkTagExists(tagId int) (bool, error) {
	var res BookTag
	query := c.db.First(&res, tagId)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if query.Error != nil {
		return false, query.Error
	}

	return true, nil
}

// UpdateBookTag 修改书籍的标签
func (c *BookTagRepo) UpdateBookTag(bookId int, tagId int) error {

	exists, err := c.checkTagExists(tagId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("tag not exists")
	}

	exists, err = DbOp.BookRepo.CheckBooId(bookId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("book not exists")
	}

	var tagRel BookTagRel

	query := c.db.Model(&BookTagRel{}).Where("book_id = ?", bookId).First(&tagRel)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		// 插入新的关联记录
		c.db.Create(&BookTagRel{BookID: bookId, TagID: tagId})
		return nil
	}
	if query.Error != nil {
		return query.Error
	}

	if tagRel.TagID == tagId {
		return nil
	}

	// 修改已有的关联记录
	tagRel.TagID = tagId
	c.db.Save(&tagRel)

	return nil
}
