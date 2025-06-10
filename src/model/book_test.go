package model

import (
	"fmt"
	"testing"
)

func TestSearchBook(t *testing.T) {

	res, err := DbOp.BookRepo.SearchBook("叔本华", 1, 10, 0)
	if err != nil {
		t.Error(err)
	}
	for _, v := range res {
		t.Log(fmt.Sprintf("%+v", v))
	}
}

func TestSearchBookTag(t *testing.T) {
	res, err := DbOp.BookRepo.SearchBook("", 2, 10, 0)
	if err != nil {
		t.Error(err)
	}
	for _, v := range res {
		t.Log(fmt.Sprintf("%+v", v))
	}
}

func TestSearchBookSearchAndTag(t *testing.T) {
	res, err := DbOp.BookRepo.SearchBook("叔本华", 3, 10, 0)

	if err != nil {
		t.Error(err)
	}

	for _, v := range res {
		t.Log(fmt.Sprintf("%+v", v))
	}

}
