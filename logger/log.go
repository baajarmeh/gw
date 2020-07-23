package logger

import (
	"fmt"
	"strings"
	"time"
)

type LogLevel uint

const (
	ERROR LogLevel = iota
	WARN
	DEBUG
	INFO
)

var logPrefix = "GW"
var logLevel = INFO
var logFormatter = "\n[$prefix] - $level,$time - $msg\n"

func SetLogPrefix(prefix string) {
	logPrefix = prefix
}

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func SetLogFormatter(formatter string) {
	logFormatter = formatter
}

func formatLog(level, format string, a ...interface{}) string {
	s := logFormatter
	t := time.Now().Format(time.RFC3339)
	s = strings.Replace(s, "$prefix", logPrefix, 1)
	s = strings.Replace(s, "$level", level, 1)
	s = strings.Replace(s, "$time", t, 1)
	msg := fmt.Sprintf(format, a...)
	return strings.Replace(s, "$msg", msg, 1)
}

func Info(format string, a ...interface{}) {
	if logLevel >= INFO {
		fmt.Printf(formatLog("INFO", format, a...))
	}
}

func Error(format string, a ...interface{}) {
	if logLevel >= ERROR {
		fmt.Printf(formatLog("ERROR", format, a...))
	}
}

func Warn(format string, a ...interface{}) {
	if logLevel >= WARN {
		fmt.Printf(formatLog("WARN", format, a...))
	}
}

func Debug(format string, a ...interface{}) {
	if logLevel >= DEBUG {
		fmt.Printf(formatLog("DEBUG", format, a...))
	}
}
