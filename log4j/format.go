package log4j

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

var (
	// 复用bytesBuffer, 避免频繁分配内存
	bytesBufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 64))
		},
	}
)

func formatLogRecord(format string, rec *logRecord) string {
	if rec == nil {
		return "\n"
	}
	if len(format) == 0 {
		return ""
	}

	// out := bytes.NewBuffer(make([]byte, 0, 64))
	out := bytesBufferPool.Get().(*bytes.Buffer)
	out.Reset()
	defer bytesBufferPool.Put(out)

	// Split the string into pieces by % signs
	pieces := bytes.Split([]byte(format), []byte{'%'})

	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'T':
				zone, _ := rec.Created.Zone()
				out.WriteString(rec.Created.Format("15:04:05"))
				out.WriteString(fmt.Sprintf(".%03d", rec.Created.UnixNano()/1e6%1000))
				out.WriteByte(' ')
				out.WriteString(zone)
			case 't':
				out.WriteString(rec.Created.Format("15:04"))
			case 'D':
				out.WriteString(rec.Created.Format("2006/01/02"))
			case 'd':
				out.WriteString(rec.Created.Format("01/02/"))
				out.WriteString(rec.Created.Format("2006")[2:])
			case 'L':
				out.WriteString(levelStrings[rec.Level])
			case 'S':
				out.WriteString(rec.Source)
			case 's':
				sources := strings.Split(rec.Source, "/")
				out.WriteString(sources[len(sources)-1])
			case 'f':
				sources := strings.Split(rec.Source, "/")
				names := strings.Split(sources[len(sources)-1], ".")
				out.WriteString(names[len(names)-1])
			case 'M':
				out.WriteString(rec.Message)
			case 'B':
				out.WriteByte('\n')
			}
			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}
	out.WriteByte('\n')

	return out.String()
}
