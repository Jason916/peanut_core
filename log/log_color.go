//jasonxu
package log

import (
	"log"
	"fmt"
)

type ColorType string

const (
	Black        ColorType = "\x1b[30m"
	Red          ColorType = "\x1b[31m"
	Green        ColorType = "\x1b[32m"
	Yellow       ColorType = "\x1b[33m"
	Blue         ColorType = "\x1b[34m"
	Magenta      ColorType = "\x1b[35m"
	Cyan         ColorType = "\x1b[36m"
	White        ColorType = "\x1b[37m"
	DefaultColor ColorType = "\x1b[39m"
)

var defcolor ColorType = DefaultColor

func DefColor() ColorType {
	return defcolor
}

func SetDefaultColor(color ColorType) {
	defcolor = color
}
func Color(color ColorType, tag string) string {
	return fmt.Sprintf("%s[%-5s]: %s", color, tag, DefColor())
}

func Info(info string, args ...interface{}) {
	info = Color(Blue, getLevelTag(LevelInfo)) + info
	log.Printf(info, args...)
}

func Trace(info string, args ...interface{}) {
	info = Color(White, getLevelTag(LevelTrace)) + info
	log.Printf(info, args...)
}

func Error(info string, args ...interface{}) {
	info = Color(Red, getLevelTag(LevelError)) + info
	log.Fatalf(info, args...)
}

func Warning(info string, args ...interface{}) {
	info = Color(Magenta, getLevelTag(LevelWarn)) + info
	log.Printf(info, args...)
}

func Success(info string, args ...interface{}) {
	info = Color(Green, getLevelTag(LevelSuccess)) + info
	log.Printf(info, args...)
}