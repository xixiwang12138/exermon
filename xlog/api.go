package xlog

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func WithContext(ctx context.Context) *Logger {
	v := ctx.Value(XLOG)
	if v == nil {
		return &Logger{
			BaseLogger: LoggerComponent,
			traceId:    "default",
		}
	}
	return v.(*Logger)
}

func NewLogger(traceID string) *Logger {
	return &Logger{
		BaseLogger: LoggerComponent,
		traceId:    traceID,
	}
}

type Logger struct {
	*BaseLogger
	traceId string
}

func (l *Logger) Write(p []byte) (n int, err error) {
	err = l.write(p, time.Now())
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

func (l *Logger) Warn(format string, args ...interface{}) {
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

func (l *Logger) LogStack(level LogLevel) {
	l.Write(debug.Stack())
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.getLevel() {
		return
	}
	now := time.Now()
	file, line := l.getCaller()
	logStr := l.getPrefix(level, file, strconv.Itoa(line))
	logStr += fmt.Sprintf(format, args...)
	_ = l.write([]byte(logStr), now)
}

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
