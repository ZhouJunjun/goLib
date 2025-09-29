/**
 * @author junjunzhou
 * @date 2022/3/10
 */
package logBuffer

import (
    "fmt"
    "github.com/ZhouJunjun/goLib/log4j"
    "github.com/ZhouJunjun/goLib/util"
    "strconv"
)

type Buffer struct {
    buffer      []byte
    total, use  int
    logLevel    log4j.Level
    printStack  bool
    runtimeSkip *int
    flagMap     map[string]bool // 自定义标志
}

func (p *Buffer) String() string {
    if p == nil {
        return ""
    }
    return string(p.buffer[:p.use])
}

func (p *Buffer) Bytes() []byte {
    if p == nil {
        return nil
    }
    return p.buffer[:p.use]
}

func (p *Buffer) ToError() error {
    if p != nil && p.use > 0 {
        return errorBuffer(p.buffer[:p.use])
    }
    return nil
}

func (p *Buffer) Set(key string, flag bool) *Buffer {
    if p == nil {
        return nil
    }

    if p.flagMap == nil {
        p.flagMap = map[string]bool{}
    }
    p.flagMap[key] = flag
    return p
}

func (p *Buffer) Is(key string) bool {
    return p != nil && p.flagMap != nil && p.flagMap[key]
}

func (p *Buffer) SetLogLevel(level log4j.Level) *Buffer {
    if p == nil {
        return nil
    }
    p.logLevel = level
    return p
}

func (p *Buffer) GetLogLevel() log4j.Level {
    if p == nil || p.logLevel == 0 {
        return log4j.INFO
    }
    return p.logLevel
}

// 废弃; 请用SetLogLevel
func (p *Buffer) SetError(isError bool) *Buffer {
    if p == nil {
        return nil
    }
    if isError {
        p.logLevel = log4j.ERROR
    } else {
        p.logLevel = log4j.INFO
    }
    return p
}

func (p *Buffer) IsError() bool {
    return p != nil && p.logLevel == log4j.ERROR
}

func (p *Buffer) SetPrintStack(printStack bool) *Buffer {
    if p == nil {
        return nil
    }
    p.printStack = printStack
    return p
}

func (p *Buffer) PrintStack() bool {
    if p == nil {
        return false
    }
    return p.printStack
}

func (p *Buffer) SetRuntimeSkip(skip int) *Buffer {
    if p == nil {
        return nil
    }
    p.runtimeSkip = &skip
    return p
}

func (p *Buffer) RuntimeSkip(defaultSkip int) int {
    if p == nil || p.runtimeSkip == nil {
        return defaultSkip
    }
    return *p.runtimeSkip
}

func (p *Buffer) Len() int {
    if p == nil {
        return 0
    }
    return p.use
}

func (p *Buffer) Wrap() *Buffer {
    if p == nil {
        return nil
    }
    return p.append([]byte("\n"))
}

func (p *Buffer) Separator(s string) *Buffer {
    if p == nil {
        return nil
    }
    if p.Len() > 0 {
        return p.append([]byte(s))
    }
    return p
}

func (p *Buffer) expand(num int) {
    if need := p.use + num - p.total; need > 0 {
        expand := util.IfInt(need%step == 0, need/step, need/step+1) * step
        newBuffer := make([]byte, p.total+expand)
        copy(newBuffer[:p.use], p.buffer)
        p.buffer, p.total = newBuffer, p.total+expand
    }
}

func (p *Buffer) append(bs []byte) *Buffer {
    if l := len(bs); l > 0 {
        p.expand(l)
        copy(p.buffer[p.use:p.use+l], bs)
        p.use += l
    }
    return p
}

func (p *Buffer) Append(i interface{}) *Buffer {
    if p == nil {
        return nil
    }
    switch i.(type) {
    case []byte:
        return p.append(i.([]byte))
    case int:
        return p.append([]byte(strconv.FormatInt(int64(i.(int)), 10)))
    case int64:
        return p.append([]byte(strconv.FormatInt(i.(int64), 10)))
    case bool:
        return p.append([]byte(util.IfString(i.(bool), "true", "false")))
    case string:
        return p.append([]byte(i.(string)))
    case error:
        if err := i.(error); err != nil {
            return p.append([]byte(err.Error()))
        } else {
            return p.append([]byte("nil"))
        }
    case *Buffer:
        if buf := i.(*Buffer); buf != nil {
            if buf.GetLogLevel() > p.GetLogLevel() {
                p.SetLogLevel(buf.GetLogLevel())
            }
            p.append(buf.Bytes())
        }
        return p
    default:
        // fmt.Println(reflect.TypeOf(i).Kind().String())
        return p.append([]byte(fmt.Sprintf("%v", i)))
    }
}

func (p *Buffer) AppendString(s string) *Buffer {
    if p == nil {
        return nil
    }
    return p.append([]byte(s))
}

func (p *Buffer) AppendInt(i int) *Buffer {
    if p == nil {
        return nil
    }
    return p.append([]byte(strconv.FormatInt(int64(i), 10)))
}

func (p *Buffer) AppendInt64(i int64) *Buffer {
    if p == nil {
        return nil
    }
    return p.append([]byte(strconv.FormatInt(int64(i), 10)))
}

func (p *Buffer) AppendBool(f bool) *Buffer {
    if p == nil {
        return nil
    }
    return p.append([]byte(util.IfString(f, "true", "false")))
}

func (p *Buffer) AppendJson(i interface{}) *Buffer {
    if p == nil {
        return nil
    }
    bs, _ := myJson.Marshal(i)
    return p.append(bs)
}

func (p *Buffer) AppendFormat(s string, a ...interface{}) *Buffer {
    if p == nil {
        return nil
    }
    return p.append([]byte(fmt.Sprintf(s, a...)))
}

func (p *Buffer) AppendMs(ms int64) *Buffer {
    if ms < 1000 {
        p.AppendFormat("%dms", ms)
    } else if ms == 1000 {
        p.AppendString("1s")
    } else {
        p.AppendFormat("%.3fs", float64(ms)/1000)
    }
    return p
}

func (p *Buffer) Insert(offset int, i interface{}) *Buffer {
    if p == nil {
        return nil
    }
    switch i.(type) {
    case []byte:
        return p.insert(offset, i.([]byte))
    case int:
        return p.insert(offset, []byte(strconv.FormatInt(int64(i.(int)), 10)))
    case int64:
        return p.insert(offset, []byte(strconv.FormatInt(i.(int64), 10)))
    case bool:
        return p.insert(offset, []byte(util.IfString(i.(bool), "true", "false")))
    case string:
        return p.insert(offset, []byte(i.(string)))
    case error:
        if err := i.(error); err != nil {
            return p.insert(offset, []byte(err.Error()))
        } else {
            return p.insert(offset, []byte("nil"))
        }
    default:
        return p.insert(offset, []byte(fmt.Sprintf("%v", i)))
    }
}

func (p *Buffer) insert(offset int, bs []byte) *Buffer {
    if offset < 0 || offset > p.use {
        panic("offset<=0 or offset > len")
    }

    if l := len(bs); l > 0 {
        newBuffer := make([]byte, p.total+l)
        if p.buffer != nil {
            copy(newBuffer[:offset], p.buffer[:offset])
            copy(newBuffer[offset+l:], p.buffer[offset:])
        }
        copy(newBuffer[offset:offset+l], bs)
        p.buffer, p.use, p.total = newBuffer, p.use+l, p.total+l
    }
    return p
}

func (p *Buffer) InsertFormat(offset int, s string, a ...interface{}) *Buffer {
    if p == nil {
        return nil
    }
    return p.insert(offset, []byte(fmt.Sprintf(s, a...)))
}

// offset >=0为下标；-1=倒数第1 -2=倒数第2...
func (p *Buffer) Del(offset, size int) *Buffer {
    if p == nil {
        return nil
    }

    if offset < 0 {
        offset = p.use + offset
    }

    if offset < 0 || offset+size > p.use {
        panic("offset incorrect or offset+size > len")
    }

    for i := offset; i+size < p.use; i++ {
        p.buffer[i] = p.buffer[i+size]
    }
    p.use -= size
    return p
}

func (p *Buffer) GetByte(i int) byte {
    if i >= 0 && i < p.Len() {
        return p.buffer[i]
    }
    return 0
}
