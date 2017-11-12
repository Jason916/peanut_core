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
func Color(color ColorType, des string) string {
	return fmt.Sprintf("%s%s%s", color, des, DefColor())
}

func Info(format string, a ...interface{}) {
	format = Color(Blue, "[INFO]: ") + format
	log.Printf(format, a...)
}

func Trace(format string, a ...interface{}) {
	format = Color(White, "[TRACE]: ") + format
	log.Printf(format, a...)
}

func Error(format string, a ...interface{}) {
	format = Color(Red, "[ERROR]: ") + format
	log.Fatalf(format, a...)
}

func Warning(format string, a ...interface{}) {
	format = Color(Magenta, "[WARN]: ") + format
	log.Fatalf(format, a...)
}

func Success(format string, a ...interface{}) {
	format = Color(Green, "[SUCC]: ") + format
	log.Fatalf(format, a...)
}
