package config

import (
	"flag"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Log(os.Args)
	logPath := flag.String("log", "reader.log", "日志文件路径")
	t.Log(*logPath)
}
