package xlog

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

type BaseLogger struct {
	level LogLevel       //日志等级
	file  io.WriteCloser //文件描述符
	mutex sync.Mutex     //互斥锁
}

func (b *BaseLogger) switchLevel(level LogLevel) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.level = level
}

func (b *BaseLogger) getLevel() LogLevel {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.level
}

func newBaseLogger(level LogLevel) (*BaseLogger, error) {
	logger := &BaseLogger{
		level: level,
	}
	return logger, nil
}

func (b *BaseLogger) write(log []byte, logTime time.Time) error {
	_, err := b.file.Write(append(log, '\n'))
	if err != nil {
		fmt.Println("Failed to write log", err)
		return err
	}
	return err
}

func (b *BaseLogger) getCaller() (string, int) {
	_, file, line, _ := runtime.Caller(3)
	return file, line
}
