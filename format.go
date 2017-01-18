package elog

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

const (
	lf         byte = 0x0A // 换行
	coreFormat      = "[%s] %v"
	timeLayout      = "2006-01-02 15:04:05.999"
)

var (
	bufPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

func (e *ELog) Debug(format string, params ...interface{}) {

}

// core func
func (e *ELog) log(msg string) {
	var buffer *bytes.Buffer
	buffer = bufPool.Get().(*bytes.Buffer)
	buffer.Reset() // 不能保证buffer是否被GC
	defer bufPool.Put(buffer)

}

func (e *ELog) logPrefix() string {
	now := time.Now()
	prefix := ""

	if len(e.cfg.TimeLayout) <= 0 {
		prefix = fmt.Sprintf(coreFormat, e.cfg.LogLevel.String(), now.Format(timeLayout))
	} else {
		prefix = fmt.Sprintf(coreFormat, e.cfg.LogLevel.String(), now.Format(e.cfg.TimeLayout))
	}

	return prefix
}
