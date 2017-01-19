package elog

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	lf          byte = 0x0A // 换行
	space       byte = 0x20 // 空格
	coreFormat       = "[%s] %v"
	timeLayout       = "2006-01-02 15:04:05.999"
	callerDepth      = 4
)

var (
	bufPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	msgPool = &sync.Pool{
		New: func() interface{} {
			return &logMessage{}
		},
	}

	errLineNo = errors.New("ELog: get lineno encounter a error.")
)

func (e *ELog) Debug(format string, params ...interface{}) {
	logMsg := msgPool.Get().(logMessage)

	logMsg.msg = fmt.Sprintf(format, params...)
	logMsg.lvl = DebugLvl

	go func() {
		e.logChan <- logMsg
	}()
}

// core func
func (e *ELog) log(logMsg logMessage) {
	var buffer *bytes.Buffer
	buffer = bufPool.Get().(*bytes.Buffer)
	buffer.Reset() // 不能保证buffer是否被GC
	defer bufPool.Put(buffer)

	e.logger.Lock()
	defer e.logger.Unlock()

	// 前缀
	buffer.WriteString(e.logPrefix(logMsg.lvl))
	buffer.WriteByte(space)

	// 行号
	if e.cfg.ShowLineNumber {
		lineNo, err := getLineNo()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			buffer.WriteString(lineNo)
			buffer.WriteByte(space)
		}
	}

	// 日志信息
	buffer.WriteString(logMsg.msg)
	buffer.WriteByte(lf)

	// 写入文件
	e.logger.f.Write(buffer.Bytes())

	// 写入标注输出
	if e.cfg.EnabledStdout {
		e.stdout.Write(buffer.Bytes())
	}
}

func (e *ELog) logPrefix(lvl LogLevel) string {
	now := time.Now()
	prefix := ""

	if len(e.cfg.TimeLayout) <= 0 {
		prefix = fmt.Sprintf(coreFormat, lvl.String(), now.Format(timeLayout))
	} else {
		prefix = fmt.Sprintf(coreFormat, lvl.String(), now.Format(e.cfg.TimeLayout))
	}

	return prefix
}

func getLineNo() (string, error) {
	_, filePath, lineNo, ok := runtime.Caller(callerDepth)
	if !ok {
		return "", errLineNo
	}

	return fmt.Sprintf("%s:%d.", filePath, lineNo), nil
}
