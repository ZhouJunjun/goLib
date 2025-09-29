/**
 * @author junjunzhou
 * @date 2024/8/2
 */
package logBuffer

import (
	"github.com/json-iterator/go"
)

const step = 16

var myJson = jsoniter.ConfigCompatibleWithStandardLibrary

func NewBuffer() *Buffer {
	return &Buffer{}
}

func NewBufferString(s string) *Buffer {
	return NewBuffer().AppendString(s)
}

func NewBufferFormat(s string, a ...interface{}) *Buffer {
	return NewBuffer().AppendFormat(s, a...)
}

type errorBuffer []byte

func (p errorBuffer) Error() string {
	return string(p)
}
