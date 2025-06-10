package util

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func IsDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func IsNoFoundError(error error) bool {
	return errors.Is(error, gorm.ErrRecordNotFound)
}
