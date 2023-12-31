package log4j

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

type ConsoleLogWriter struct {
	rec    chan *logRecord
	format string
	level  level
}

func (p *ConsoleLogWriter) SetFormat(format string) {
	p.format = format
}

func NewConsoleLogWriter(level level) *ConsoleLogWriter {
	writer := &ConsoleLogWriter{
		rec:    make(chan *logRecord, LogBufferLength),
		format: "[%D %T] [%L] (%S) %M",
		level:  level,
	}
	go writer.run(stdout)
	return writer
}

func (p *ConsoleLogWriter) run(out io.Writer) {
	for rec := range p.rec {
		_, _ = fmt.Fprint(out, formatLogRecord(p.format, rec))
	}
}

func (p *ConsoleLogWriter) LogWrite(rec *logRecord) {
	p.rec <- rec
}

func (p *ConsoleLogWriter) Close() {
	close(p.rec)
}

func (p *ConsoleLogWriter) IsPrivate() bool {
	return false
}

func (p *ConsoleLogWriter) GetLevel() level {
	return p.level
}
