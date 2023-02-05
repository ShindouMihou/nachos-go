package nachos

import "log"

type LoggingLevel int32

const (
	Trace LoggingLevel = 0
	Info  LoggingLevel = 1
	Warn  LoggingLevel = 2
	Error LoggingLevel = 3
)

var Logger = func(level LoggingLevel, v ...any) {
	log.Println(v)
}
