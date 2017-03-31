// Copyright (C) 2017, No.20 <zdt3476@163.com>. All rights reserved

package elog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	lf          byte = 0x0A // 换行
	space       byte = 0x20 // 空格
	separator        = string(filepath.Separator)
	coreFormat       = "[%s] %v"
	timeLayout       = "2006-01-02 15:04:05.999"
	callerDepth      = 3
)

var (
	errLineNo = errors.New("[ELOG]:Get lineno encounter a error.")
)

func (e *ELog) baseLog(lvl LogLevel, msg string) {
	if e.cfg.LogLevel > lvl {
		return
	}
	logMsg := getMsg()
	defer putMsg(logMsg)

	logMsg.msg = msg
	logMsg.lvl = lvl
	line, err := getLineNo(e.cfg.ShortFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get line number encounter a error.")
	}
	logMsg.lineNo = line

	e.log(logMsg)
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

func (e *ELog) Panic(format string, params ...interface{}) {
	e.baseLog(PanicLvl, fmt.Sprintf(format, params...))
}

func (e *ELog) Fatal(format string, params ...interface{}) {
	e.baseLog(FatalLvl, fmt.Sprintf(format, params...))
}

// core func
// panic stop from here
func (e *ELog) log(logMsg *logMessage) {
	buffer := getBuffer()
	defer putBuffer(buffer)

	e.logger.Lock()
	defer e.logger.Unlock()

	buffer.WriteString(e.logPrefix(logMsg.lvl))
	buffer.WriteByte(space)

	if e.cfg.ShowLineNumber {
		buffer.WriteString(logMsg.lineNo)
		buffer.WriteByte(space)
	}

	buffer.WriteString(logMsg.msg)
	buffer.WriteByte(lf)

	e.logger.f.Write(buffer.Bytes())

	if e.cfg.EnabledStdout {
		e.stdout.Write(buffer.Bytes())
	}

	if logMsg.lvl == FatalLvl {
		os.Exit(-1)
	} else if logMsg.lvl == PanicLvl {
		panic("")
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

func getLineNo(short bool) (string, error) {
	_, filePath, lineNo, ok := runtime.Caller(callerDepth)
	if !ok {
		return "", errLineNo
	}

	if short {
		idx := strings.LastIndex(filePath, separator)
		if idx != -1 {
			filePath = filePath[idx+1:]
		}
	}

	return fmt.Sprintf("%s:%d.", filePath, lineNo), nil
}
