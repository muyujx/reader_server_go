package model

import (
	"gorm.io/gorm"
	"time"
)

type FavoriteBookRepo struct {
	db *gorm.DB
}

type FavoriteBook struct {
	ID int `gorm:"primaryKey"`

	UserId int

	BookId int

	Page int

	// 上次阅读时间
	LastRead int64

	// 创建时间
	CreateTime int64

	// 修改时间
	UpdateTime int64
}

func NewFavoriteBookRepo(db *gorm.DB) *FavoriteBookRepo {
	return &FavoriteBookRepo{db: db}
}

func (*FavoriteBook) TableName() string {
	return "favorite_book"
}

func (r *FavoriteBookRepo) AddFavoriteBook(record *FavoriteBook) error {
	return r.db.Create(record).Error
}

func (r *FavoriteBookRepo) GetFavoriteBook(userId int, bookId int) (FavoriteBook, error) {
	res := FavoriteBook{}
	err := r.db.Where("user_id = ? AND book_id = ?", userId, bookId).First(&res).Error
	return res, err
}

func (r *FavoriteBookRepo) DelFavoriteBook(userId int, bookId int) error {
	return r.db.Where("user_id = ? AND book_id = ?", userId, bookId).Delete(&FavoriteBook{}).Error
}

func (r *FavoriteBookRepo) UpdateUserCurPage(userId int, bookId int, page int, lastRead int64) error {
	return r.db.Model(&FavoriteBook{}).
		Where("user_id = ? AND book_id = ?", userId, bookId).
		UpdateColumns(map[string]any{
			"page":      page,
			"last_read": lastRead,
		}).Error
}

// CountBook 用户收藏书籍总数
func (r *FavoriteBookRepo) CountBook(userId int) (int, error) {
	var res int64
	err := r.db.Model(&FavoriteBook{}).Where("user_id = ?", userId).Count(&res).Error
	return int(res), err
}

type BookFavoriteInfo struct {
	BookId      int    `json:"bookId"`
	BookName    string `json:"bookName"`
	TotalPage   int    `json:"totalPage"`
	CoverPic    string `json:"coverPic"`
	BigCoverPic string `json:"bigCoverPic"`
	TagId       int    `json:"tagId"`
	Page        int    `json:"page"`
	LastRead    int64  `json:"lastRead"`
	ReadingCost int    `json:"readingCost"`
}

func (r *FavoriteBookRepo) GetFavoriteBookList(userId int, page int, pageSize int) ([]BookFavoriteInfo, error) {
	var res []BookFavoriteInfo

	err := r.db.Table("favorite_book").
		Select(` book.id AS book_id, 
						book_name, 
						total_page, 
						cover_pic, 
						big_cover_pic, 
						COALESCE(tag_id, -1) AS tag_id, 
						favorite_book.page AS page,
						favorite_book.last_read AS last_read,
						reading_cost
		`).
		Joins(" JOIN book ON book.id = favorite_book.book_id ").
		Joins(" LEFT JOIN book_tag_rel ON book.id = book_tag_rel.book_id ").
		Where("user_id = ?", userId).
		Order("favorite_book.last_read DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&res).Error

	return res, err
}

func (r *FavoriteBookRepo) CheckFavoriteBook(userId int, bookIds []int) (map[int]interface{}, error) {
	var list []int
	err := r.db.Table("favorite_book").
		Select("book_id").
		Where("user_id = ? AND book_id IN (?)", userId, bookIds).
		Find(&list).Error

	if err != nil {
		return nil, err
	}

	var resMap = make(map[int]interface{})
	for _, id := range list {
		resMap[id] = nil
	}
	return resMap, nil
}

// AddBookReadingTime 增加书籍阅读时间
func (r *FavoriteBookRepo) AddBookReadingTime(userId int, bookId int, readingTime int) error {
	return r.db.Model(&FavoriteBook{}).
		Where("book_id = ? AND user_id = ?", bookId, userId).
		UpdateColumns(map[string]any{
			"reading_cost": gorm.Expr("reading_cost + ?", readingTime),
			"update_time":  time.Now().Unix(),
		}).Error
}
