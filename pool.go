// Copyright (C) 2017, No.20 <zdt3476@163.com>. All rights reserved

package elog

import (
	"bytes"
	"sync"
)

var (
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	msgPool = &sync.Pool{
		New: func() interface{} {
			return &logMessage{}
		},
	}
)

func getBuffer() (buf *bytes.Buffer) {
	return bufferPool.Get().(*bytes.Buffer)
}

// The buffer is reset before it is put back into circulation.
func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

func getMsg() (msg *logMessage) {
	return msgPool.Get().(*logMessage)
}

func putMsg(msg *logMessage) {
	msgPool.Put(msg)
}
