package model

import (
	"errors"
	"gorm.io/gorm"
	"testing"
)

func TestGetReadHistory(t *testing.T) {
	res, err := DbOp.ReadHistoryRepo.GetReadHistory(0, 0, "")
	t.Log("GetReadHistory res:", res, "err:", err)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fail()
	}
}

func TestListReadHistory(t *testing.T) {
	res, err := DbOp.ReadHistoryRepo.ListReadHistory(10, "2025-10-13", "2025-10-15")
	if err != nil {
		t.Fail()
	}
	t.Log(res)
}
