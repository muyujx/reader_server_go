package model

import (
	"testing"
	"time"
)

func Test_AddFavoriteBook(t *testing.T) {
	curTime := time.Now().Unix()
	err := DbOp.FavoriteBookRepo.AddFavoriteBook(&FavoriteBook{
		BookId:     1,
		UserId:     1,
		CreateTime: curTime,
		UpdateTime: curTime,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestFavoriteBookRepo_GetFavoriteBook(t *testing.T) {
	res, err := DbOp.FavoriteBookRepo.GetFavoriteBook(1, 1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", res)
}

func TestFavoriteBookRepo_DelFavoriteBook(t *testing.T) {
	err := DbOp.FavoriteBookRepo.DelFavoriteBook(1, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestFavoriteBookRepo_CountBook(t *testing.T) {
	res, err := DbOp.FavoriteBookRepo.CountBook(1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", res)
}
