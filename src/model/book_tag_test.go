package model

import (
	"testing"
)

func TestGetAllTag(t *testing.T) {
	tags, err := DbOp.BookTagRepo.GetAllTag()
	if err != nil {
		t.Error(err)
	}
	t.Log(tags)
}
