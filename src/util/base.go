package util

import (
	"muyu.com/reader_server_go/v1/src/logger"
)

func GoSafe(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Warning("GoSafe recover from panic", r)
			}
		}()
		fn()
	}()
}
