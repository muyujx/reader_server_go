package logger

import (
	"fmt"
	"github.com/fatih/color"
	"muyu.com/reader_server_go/v1/src/config"
	"sync"
	"time"
)

const (
	// LevelError 错误
	LevelError = iota
	// LevelWarning 警告
	LevelWarning
	// LevelInfo 提示
	LevelInfo
	// LevelDebug 除错
	LevelDebug
)

// Logger 日志
type Logger struct {
	level int
	mu    sync.Mutex
}

// 日志颜色
var colors = map[string]func(a ...interface{}) string{
	"Warning": color.New(color.FgYellow).Add(color.Bold).SprintFunc(),
	"Panic":   color.New(color.BgRed).Add(color.Bold).SprintFunc(),
	"Error":   color.New(color.FgRed).Add(color.Bold).SprintFunc(),
	"Info":    color.New(color.FgCyan).Add(color.Bold).SprintFunc(),
	"Debug":   color.New(color.FgWhite).Add(color.Bold).SprintFunc(),
}

// 不同级别前缀与时间的间隔，保持宽度一致
var spaces = map[string]string{
	"Warning": "",
	"Panic":   "  ",
	"Error":   "  ",
	"Info":    "   ",
	"Debug":   "  ",
}

// Println 打印
func (ll *Logger) Println(prefix string, msg string) {
	c := color.New()

	ll.mu.Lock()
	defer ll.mu.Unlock()

	_, _ = c.Printf(
		"%s%s %s %s\n",
		colors[prefix]("["+prefix+"]"),
		spaces[prefix],
		time.Now().Format("2006-01-02 15:04:05"),
		msg,
	)
}

// Panic 极端错误
func Panic(format string, v ...any) {
	ll := log()
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	ll.Println("Panic", msg)
	panic(msg)
}

// Error 错误
func Error(format string, v ...any) {
	ll := log()
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	ll.Println("Error", msg)
}

// Warning 警告
func Warning(format string, v ...any) {
	ll := log()
	if LevelWarning > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	ll.Println("Warning", msg)
}

// Info 信息
func Info(format string, v ...any) {
	ll := log()
	if LevelInfo > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	ll.Println("Info", msg)
}

// Debug 校验
func Debug(format string, v ...any) {
	ll := log()
	if LevelDebug > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	ll.Println("Debug", msg)
}

var GlobalLogger *Logger

// Log 返回日志对象
func log() *Logger {
	if GlobalLogger == nil {
		logLevel := config.Config.LogLevel

		level := LevelDebug
		switch logLevel {
		case "info":
			level = LevelInfo
		case "warn":
			level = LevelWarning
		case "error":
			level = LevelError
		default:
			level = LevelDebug
		}

		GlobalLogger = &Logger{
			level: level,
		}
	}
	return GlobalLogger
}
