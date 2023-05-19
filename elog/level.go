package elog

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var logLevelNames = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

//var logLevelColors = map[LogLevel]string{
//	DEBUG:   "\033[34m", // Blue
//	INFO:    "\033[32m", // Green
//	WARNING: "\033[33m", // Yellow
//	ERROR:   "\033[31m", // Red
//	FATAL:   "\033[35m", // Magenta
//}

type Logger struct {
	*BaseLogger
	traceId string
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.getLevel() {
		return
	}
	now := time.Now()
	file, line := l.getCaller()
	logStr := fmt.Sprintf("[%s][%s][%s] %s:%d  ", l.traceId, now.Format("2006-01-02 15:04:05"), logLevelNames[l.getLevel()], file, line)
	logStr += fmt.Sprintf(format, args...)
	l.writeToFile(logStr, now)
}
