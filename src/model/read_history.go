package model

import (
	"errors"
	"gorm.io/gorm"
)

// ReadHistory 每天的阅读历史记录
type ReadHistory struct {
	ID int `gorm:"primaryKey"`

	UserId int

	BookId int

	// 时间字符串 2025-10-12
	DayStr string

	// 当前的读书总时间 单位 s
	ReadingCost int

	// 当天读书起始页
	StartPage int

	// 当天读书结束页
	EndPage int

	// 创建时间
	CreateTime int64

	// 修改时间
	UpdateTime int64
}

func (*ReadHistory) TableName() string {
	return "read_history"
}

type ReadHistoryRepo struct {
	db *gorm.DB
}

func NewReadHistoryRepo(db *gorm.DB) *ReadHistoryRepo {
	return &ReadHistoryRepo{db: db}
}

// GetReadHistory 获取阅读历史记录
func (rp *ReadHistoryRepo) GetReadHistory(
	userId int,
	bookId int,
	dayStr string) (*ReadHistory, error) {
	var res ReadHistory
	err := rp.db.Where("user_id = ? and book_id = ? and day_str = ?",
		userId, bookId, dayStr).First(&res).Error
	return &res, err
}

// CheckExists 检查历史记录是否存在
func (rp *ReadHistoryRepo) CheckExists(
	userId int,
	bookId int,
	dayStr string) (bool, error) {
	_, err := rp.GetReadHistory(userId, bookId, dayStr)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// AddReadHistory 创建阅读历史记录
func (rp *ReadHistoryRepo) AddReadHistory(record *ReadHistory) error {
	return rp.db.Create(record).Error
}

// UpdateReadingPageCost 更新当天阅读时间和阅读页数
func (rp *ReadHistoryRepo) UpdateReadingPageCost(
	userId int,
	bookId int,
	dayStr string,
	readingTime int,
	page int) error {
	return rp.db.Model(&ReadHistory{}).
		Where("user_id = ? and book_id = ? and day_str = ?", userId, bookId, dayStr).
		UpdateColumns(map[string]any{
			"reading_cost": gorm.Expr("reading_cost + ?", readingTime),
			"end_page":     gorm.Expr("GREATEST(start_page, ?)", page),
		}).Error
}

type ListReadHistoryDto struct {
	ReadHistory

	BookName string

	CoverPic string

	BigCoverPic string
}

func (rp *ReadHistoryRepo) ListReadHistory(
	userId int,
	startDateStr,
	endDateStr string) ([]ListReadHistoryDto, error) {
	var res []ListReadHistoryDto
	err := rp.db.Model(&ReadHistory{}).
		Select(" read_history.*, book_name, cover_pic, big_cover_pic ").
		Joins("JOIN book ON book.id = read_history.book_id").
		Where(" user_id = ? AND day_str >= ? AND day_str <= ?", userId, startDateStr, endDateStr).
		Order("read_history.day_str, book.id").
		Find(&res).Error
	return res, err
}
