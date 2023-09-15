/**
 * @author junjunzhou
 * @date 2022/3/10
 */
package logBuffer

import (
	"fmt"
	"strconv"
	"unsafe"
)

func NewBufferString(s string) *Buffer {
	return &Buffer{buf: []byte(s)}
}

func NewBufferFormat(s string, a ...interface{}) *Buffer {
	return &Buffer{buf: []byte(fmt.Sprintf(s, a...))}
}

type Buffer struct {
	buf         []byte
	isError     bool
	printStack  bool
	runtimeSkip *int
	tag         string
}

func (p *Buffer) String() string {
	return *(*string)(unsafe.Pointer(&p.buf))
}

func (p *Buffer) SetError(isError bool) *Buffer {
	p.isError = isError
	return p
}

func (p *Buffer) IsError() bool {
	return p.isError
}

func (p *Buffer) SetPrintStack(printStack bool) *Buffer {
	p.printStack = printStack
	return p
}

func (p *Buffer) PrintStack() bool {
	return p.printStack
}

func (p *Buffer) SetTag(tag string) *Buffer {
	p.tag = tag
	return p
}

func (p *Buffer) Tag() string {
	return p.tag
}

func (p *Buffer) SetRuntimeSkip(skip int) *Buffer {
	p.runtimeSkip = &skip
	return p
}

func (p *Buffer) RuntimeSkip(defaultSkip int) int {
	if p.runtimeSkip == nil {
		return defaultSkip
	}
	return *p.runtimeSkip
}

func (p *Buffer) Len() int {
	return len(p.buf)
}

func (p *Buffer) Append(bs []byte) *Buffer {
	p.buf = append(p.buf, bs...)
	return p
}

func (p *Buffer) AppendString(s string) *Buffer {
	p.buf = append(p.buf, s...)
	return p
}

func (p *Buffer) AppendInt(i int) *Buffer {
	p.buf = append(p.buf, strconv.FormatInt(int64(i), 10)...)
	return p
}

func (p *Buffer) AppendInt64(i int64) *Buffer {
	p.buf = append(p.buf, strconv.FormatInt(i, 10)...)
	return p
}

func (p *Buffer) AppendFormat(s string, a ...interface{}) *Buffer {
	tmpS := fmt.Sprintf(s, a...)
	p.buf = append(p.buf, tmpS...)
	return p
}

func (p *Buffer) Insert(offset int, bs []byte) *Buffer {
	if offset < 0 || offset > len(p.buf) {
		panic("")
	}

	if bs == nil {
		bs = []byte("null")
	}

	if size := len(bs); size > 0 {
		newBuf := make([]byte, len(p.buf)+size)
		copy(newBuf[:offset], p.buf[:offset])
		copy(newBuf[offset:offset+size], bs)
		copy(newBuf[offset+size:], p.buf[offset:])
		p.buf = newBuf
	}
	return p
}

func (p *Buffer) InsertString(offset int, s string) *Buffer {
	p.Insert(offset, []byte(s))
	return p
}

func (p *Buffer) InsertInt(offset int, i int) *Buffer {
	p.Insert(offset, []byte(strconv.FormatInt(int64(i), 10)))
	return p
}

func (p *Buffer) InsertInt64(offset int, i int64) *Buffer {
	p.Insert(offset, []byte(strconv.FormatInt(i, 10)))
	return p
}

func (p *Buffer) InsertFormat(offset int, s string, a ...interface{}) *Buffer {
	tmpS := fmt.Sprintf(s, a...)
	p.Insert(offset, []byte(tmpS))
	return p
}
