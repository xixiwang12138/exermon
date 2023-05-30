package elog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type BaseLogger struct {
	level      LogLevel     //日志等级
	logDir     string       //输入的日志文件夹路径
	file       *os.File     //文件描述符
	currentDay atomic.Int32 // 20220110, use atomic

	mutex sync.Mutex //互斥锁
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

func (b *BaseLogger) geFile() *os.File {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.file
}

func newBaseLogger(level LogLevel, logDir string) (*BaseLogger, error) {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}
	logger := &BaseLogger{
		level:  level,
		logDir: logDir,
	}
	if err := logger.lockedOpenFile(time.Now()); err != nil {
		return nil, err
	}
	return logger, nil
}

func (b *BaseLogger) lockedOpenFile(logTime time.Time) error {
	fileName := logTime.Format("2006.01.02") + ".log"
	filePath := filepath.Join(b.logDir, fileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.currentDay.Store(dayIntFlag(logTime))
	b.file = file
	return nil
}

func (b *BaseLogger) writeToFile(logStr string, logTime time.Time) error {
	b.tryOpenNewFile(logTime)
	_, err := b.geFile().WriteString(logStr + "\n")
	if err != nil {
		fmt.Println("Failed to write log", err)
	}
	return err
}

func dayIntFlag(t time.Time) int32 {
	y, m, d := t.Date()
	return int32(y*10000 + int(m)*100 + d)
}

func (b *BaseLogger) tryOpenNewFile(logTime time.Time) {
	logDay := dayIntFlag(logTime)
	if b.currentDay.Load() == logDay {
		return
	}
	// should open new file
	if err := b.lockedOpenFile(logTime); err != nil {
		log.Println("open log file error: ", err.Error())
	}
}

func (b *BaseLogger) getCaller() (string, int) {
	_, file, line, _ := runtime.Caller(3)
	return file, line
}
