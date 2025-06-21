package es

import (
	"muyu.com/reader_server_go/v1/src/config"
	"testing"
)

func init() {
	testing.Init()
	config.InitConfig()
	InitEs()
}
