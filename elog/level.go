package elog

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
	SILENT = 100
)

var logLevelToName = map[LogLevel][]byte{
	DEBUG: []byte("DEBUG"),
	INFO:  []byte("INFO"),
	WARN:  []byte("WARN"),
	ERROR: []byte("ERROR"),
	FATAL: []byte("FATAL"),
}

var logNameToLevel = map[string]LogLevel{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
	"FATAL": FATAL,
}

type Logger struct {
	*BaseLogger
	traceId string
}

func (l *Logger) Write(p []byte) (n int, err error) {
	err = l.writeToFile(string(p), time.Now())
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (l *Logger) Debug(format string, args ...interface{}) {
	format, args = replaceJsonHolder(format, args)
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	format, args = replaceJsonHolder(format, args)
	l.log(INFO, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	format, args = replaceJsonHolder(format, args)
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	format, args = replaceJsonHolder(format, args)
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	format, args = replaceJsonHolder(format, args)
	l.log(FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.getLevel() {
		return
	}
	now := time.Now()
	file, line := l.getCaller()
	logStr := l.getPrefix(level, file, strconv.Itoa(line))
	logStr += fmt.Sprintf(format, args...)
	_ = l.writeToFile(logStr, now)
}

const (
	left  byte = '['
	right byte = ']'
	black byte = ' '
	colon byte = ':'
)

var (
	rightAndLeft []byte = []byte("][")
)

func (l *Logger) getPrefix(level LogLevel, file, line string) string {
	now := time.Now()
	sb := strings.Builder{}
	sb.Grow(64)
	sb.WriteByte(left)
	sb.WriteString(l.traceId)
	sb.Write(rightAndLeft)
	sb.WriteString(now.Format("15:04:05.000"))
	sb.Write(rightAndLeft)
	sb.Write(logLevelToName[level])
	sb.WriteByte(right)
	sb.WriteByte(black)
	sb.WriteString(file)
	sb.WriteByte(colon)
	sb.WriteString(line)
	sb.WriteByte(black)
	return sb.String()
}
