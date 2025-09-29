/**
 * @author junjunzhou
 * @date 2023/1/10
 */
package log4j

import (
	"time"
)

type LogProperty struct {
	Filename string
	Format   string
	Rotate   bool
	Maxsize  int
	MaxLines int
	Daily    bool
	KeepDay  int64
	Private  bool
}

type logRecord struct {
	Level   Level     // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Source  string    // The message source
	Message string    // The log message
}

type logWriter interface {
	LogWrite(rec *logRecord)
	Close()
	IsPrivate() bool
	GetLevel() Level
}

type xmlProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlFilter struct {
	Enabled  string        `xml:"enabled,attr"`
	Tag      string        `xml:"tag"`
	Level    string        `xml:"level"`
	Type     string        `xml:"type"`
	Property []xmlProperty `xml:"property"`
}

type xmlLoggerConfig struct {
	Filter []xmlFilter `xml:"filter"`
}
