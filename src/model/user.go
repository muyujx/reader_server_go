package model

import (
	"errors"
	"gorm.io/gorm"
	"strings"
)

type UserRepo struct {
	db *gorm.DB
}

type User struct {
	ID int `gorm:"primaryKey"`

	// 账号
	Account string

	// 密码
	Password string

	// 上次登录时间
	LastLogin int64

	// 角色
	Role int

	// 创建时间
	CreateTime int64

	// 修改时间
	UpdateTime int64
}

const (

	// Role_User 普通用户
	Role_User = 0

	// Role_Admin 管理员
	Role_Admin = 1
)

func NewUseRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (*User) TableName() string {
	return "user"
}

func (c *UserRepo) GetUserByAccount(account string) (res User, err error) {
	account = strings.TrimSpace(account)
	if len(account) == 0 {
		return res, errors.New("account is empty")
	}
	err = c.db.Where("account = ?", account).First(&res).Error
	return res, err
}

func (c *UserRepo) GetAllUser() []User {
	var users []User
	c.db.Find(&users)
	return users
}

func (c *UserRepo) UpdateLoginTime(id int, lastLogin int64) (err error) {
	return c.db.Model(&User{}).Where("id = ?", id).Update("last_login", lastLogin).Error
}

func (c *UserRepo) CreateUser(user *User) error {
	return c.db.Create(user).Error
}
