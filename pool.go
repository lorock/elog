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

// GetBuffer returns a buffer from the pool.
func GetBuffer() (buf *bytes.Buffer) {
	return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the pool.
// The buffer is reset before it is put back into circulation.
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

func GetMsg() (msg *logMessage) {
	return msgPool.Get().(*logMessage)
}

func PutMsg(msg *logMessage) {
	msgPool.Put(msg)
}
