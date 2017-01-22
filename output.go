package elog

import (
	"bytes"
	"errors"
	"fmt"
	"log"
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
	callerDepth      = 3
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

func (e *ELog) baseLog(lvl LogLevel, msg string) {
	if e.cfg.LogLevel > lvl {
		return // 屏蔽打印
	}
	logMsg := msgPool.Get().(*logMessage)

	logMsg.msg = msg
	logMsg.lvl = lvl
	line, err := getLineNo()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get line number encounter a error.")
	}
	logMsg.lineNo = line

	e.waitWrite.Wait()
	e.logChan <- logMsg
}

func (e *ELog) Debug(format string, params ...interface{}) {
	e.baseLog(DebugLvl, fmt.Sprintf(format, params...))
}

func (e *ELog) Info(format string, params ...interface{}) {
	e.baseLog(InfoLvl, fmt.Sprintf(format, params...))
}

func (e *ELog) Warn(format string, params ...interface{}) {
	e.baseLog(WarnLvl, fmt.Sprintf(format, params...))
}

func (e *ELog) Error(format string, params ...interface{}) {
	e.baseLog(ErrorLvl, fmt.Sprintf(format, params...))
}

// Panic由调用者处理
func (e *ELog) Panic(format string, params ...interface{}) {
	e.baseLog(PanicLvl, fmt.Sprintf(format, params...))
}

// 这里会退出程序,调用os.Exit()
func (e *ELog) Fatal(format string, params ...interface{}) {
	e.baseLog(FatalLvl, fmt.Sprintf(format, params...))
}

// core func
func (e *ELog) log(logMsg *logMessage) {
	defer func() {
		if err := recover(); err != nil {
			// avoid panic main channel.
			log.Printf("Recover from a error.Error(%v)\n", err)
		}
	}()

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
		buffer.WriteString(logMsg.lineNo)
		buffer.WriteByte(space)
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

// 失败返回空字符串
func getLineNo() (string, error) {
	_, filePath, lineNo, ok := runtime.Caller(callerDepth)
	if !ok {
		return "", errLineNo
	}

	return fmt.Sprintf("%s:%d.", filePath, lineNo), nil
}
