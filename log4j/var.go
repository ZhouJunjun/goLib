/**
 * @author junjunzhou
 * @date 2023/1/10
 */
package log4j

import (
	"bytes"
	"sync"
)

var (
	// 日志缓冲区长度, 设置成可导出, 如果有日志多, 可改大
	LogBufferLength = 32

	defaultFormat = "[%D %T] [%L] (%S) %M"

	// 复用bytes.Buffer|[]byte, 避免频繁分配内存
	bytesBufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 64))
		},
	}
	bytesPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 64<<10)
		},
	}
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
)

var (
	levelStrings = [...]string{"DEBG", "INFO", "WARN", "EROR"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}
