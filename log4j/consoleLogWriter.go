// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4j

import (
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

type ConsoleLogWriter struct {
	rec    chan *logRecord
	format string
	level  Level
}

func (p *ConsoleLogWriter) SetFormat(format string) {
	p.format = format
}

func NewConsoleLogWriter(level Level) *ConsoleLogWriter {
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
		_, _ = fPrintFormatLog(out, p.format, rec)
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

func (p *ConsoleLogWriter) GetLevel() Level {
	return p.level
}
