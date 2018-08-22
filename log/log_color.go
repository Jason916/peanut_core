//jasonxu
package log

import (
	"log"
	"fmt"
	"bytes"
)

const (
	Black        = "\x1b[30m"
	Red          = "\x1b[31m"
	Green        = "\x1b[32m"
	Yellow       = "\x1b[33m"
	Blue         = "\x1b[34m"
	Magenta      = "\x1b[35m"
	Cyan         = "\x1b[36m"
	White        = "\x1b[37m"
	DefaultColor = "\x1b[39m"
)

var defcolor string = DefaultColor

func DefColor() string {
	return defcolor
}

func SetDefaultColor(color string) {
	defcolor = color
}

func Color(color, tag string) string {
	return fmt.Sprintf("%s[%-5s]: %s", color, tag, DefColor())
}

func Info(info string, args ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(Color(Blue, getLevelTag(LevelInfo)))
	buffer.WriteString(info)
	s := buffer.String()
	log.Printf(s, args...)
}

func Trace(info string, args ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(Color(White, getLevelTag(LevelTrace)))
	buffer.WriteString(info)
	s := buffer.String()
	log.Printf(s, args...)
}

func Error(info string, args ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(Color(Red, getLevelTag(LevelError)))
	buffer.WriteString(info)
	s := buffer.String()
	log.Fatalf(s, args...)
}

func Warning(info string, args ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(Color(Magenta, getLevelTag(LevelWarn)))
	buffer.WriteString(info)
	s := buffer.String()
	log.Printf(s, args...)
}

func Success(info string, args ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(Color(Green, getLevelTag(LevelSuccess)))
	buffer.WriteString(info)
	s := buffer.String()
	log.Printf(s, args...)
}