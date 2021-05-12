/**
 * @Author: koulei
 * @Description: TODO
 * @File:  logger
 * @Version: 1.0.0
 * @Date: 2021/5/12 21:26
 */

package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
	"time"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

const (
	info  = "[INFO]"
	warn  = "[WARN]"
	err   = "[ERROR]"
	fatal = "[FATAL]"
)

const (
	green        = "\033[97;42m"
	white        = "\033[90;47m"
	yellow       = "\033[90;43m"
	red          = "\033[97;41m"
	blue         = "\033[97;44m"
	magenta      = "\033[97;45m"
	cyan         = "\033[97;46m"
	reset        = "\033[0m"
	greenBold    = "\033[32;1m"
	yellowBold   = "\033[33;1m"
	redBold      = "\033[31;1m"
	magentaBold  = "\033[35;1m"
	greenLight   = "\033[32m"
	yellowLight  = "\033[33m"
	redLight     = "\033[31m"
	magentaLight = "\033[35m"
)

type Logger struct {
	BeginTime time.Time
	EndTime   time.Duration
	Level     uint
	ColorCode string
	Tag       string
	Message   []*msg
}

type msg struct {
	File     string
	FuncName string
	Msg      interface{}
}

func Start(c *gin.Context) {
	logger := CreateLogger()
	logger.BeginTime = time.Now()
	c.Set("logger", logger)
}

func New(c *gin.Context) *Logger {
	logger, exists := c.Get("logger")
	if !exists {
		fmt.Println("logger没有初始化，请在第一层路由中间件中初始化")
		c.Abort()
		return nil
	} else {
		return logger.(*Logger)
	}
}

func CreateLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg ...interface{}) {
	l.Level = LevelInfo
	l.Tag = info
	l.prepend(msg)
}
func (l *Logger) Warn(msg ...interface{}) {
	l.Level = LevelWarning
	l.Tag = warn
	l.prepend(msg)
}
func (l *Logger) Err(msg ...interface{}) {
	l.Level = LevelError
	l.Tag = err
	l.prepend(msg)
}

func (l *Logger) Fatal(msg ...interface{}) {
	l.Level = LevelFatal
	l.Tag = fatal
	l.prepend(msg)
	panic(msg)
}

func (l *Logger) Print() {
	fmt.Print("\n-------------------------------- Log Trace -----------------------------------\n\n")
	for _, v := range l.Message {
		fmt.Printf("%v\n%v\n%v\n\n", v.File, v.FuncName, v.Msg)
	}
	fmt.Printf("%v%v%v\n\n", yellowLight, time.Since(l.BeginTime), reset)
	fmt.Print("-------------------------------- Trace End -----------------------------------\n\n")
}

func (l *Logger) prepend(v []interface{}) {
	pc, file, line, _ := runtime.Caller(2)
	pcName := runtime.FuncForPC(pc).Name()
	var log interface{} = ""
	func() {
		for _, v := range v {
			log = fmt.Sprintf("%v%v", log, v)
		}
	}()

	switch l.Level {
	case LevelInfo:
		l.ColorCode = greenBold
	case LevelWarning:
		l.ColorCode = yellowBold
	case LevelError:
		l.ColorCode = redBold
	}

	l.Message = append(l.Message, &msg{
		File:     fmt.Sprintf("%v%s %v文件: %v:%v", l.ColorCode, l.Tag, reset, file, line),
		FuncName: fmt.Sprintf("%v%s %v函数: %v", l.ColorCode, l.Tag, reset, pcName),
		Msg:      fmt.Sprintf("%v%s %v日志: %v", l.ColorCode, l.Tag, reset, log),
	})
}
