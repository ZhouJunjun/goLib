package log4j

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type loggerHandler struct {
	logWriterMap       map[string]logWriter
	lock               sync.RWMutex
	defaultLogFilePath string
}

func newDefaultLogger(lvl Level) *loggerHandler {
	return &loggerHandler{
		lock: sync.RWMutex{},
		logWriterMap: map[string]logWriter{
			"stdout": NewConsoleLogWriter(lvl)},
	}
}

// Close all open loggers
func (p *loggerHandler) close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for name, logWriter := range p.logWriterMap {
		logWriter.Close()
		delete(p.logWriterMap, name)
	}
}

func (p *loggerHandler) closeByTag(logTag string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if logWriter, ok := p.logWriterMap[logTag]; ok {
		logWriter.Close()
		delete(p.logWriterMap, logTag)
	}
}

func (p *loggerHandler) getLogFilePath() string {
	return p.defaultLogFilePath
}

func (p *loggerHandler) getLogWriterMap() map[string]logWriter {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.logWriterMap
}

func (p *loggerHandler) addLogString(runtimeSkip int, lvl Level, withStack bool, tag string, format string, args ...interface{}) {

	src, msg := getSrcAndMsg(runtimeSkip, withStack, format, args...)

	rec := &logRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: msg,
	}

	p.print(rec, lvl, tag)
}

func (p *loggerHandler) print(rec *logRecord, lvl Level, tag string) {
	logWriterMap := p.getLogWriterMap()

	// 指定tag 且 对应的tag文件存在且私有, 只写私有
	if tag != "" {
		for tagName, logWriter := range logWriterMap {
			if tagName == tag && logWriter.IsPrivate() {
				if lvl >= logWriter.GetLevel() {
					logWriter.LogWrite(rec)
				}
				return
			}
		}
	}

	isPrinted := false
	for _, logWriter := range logWriterMap {
		if lvl >= logWriter.GetLevel() && !logWriter.IsPrivate() {
			logWriter.LogWrite(rec)
			isPrinted = true
		}
	}

	/*for tagName, logWriter := range logWriterMap {
	    if lvl >= logWriter.GetLevel() {
	        if tag == "" && !logWriter.IsPrivate() {
	            logWriter.LogWrite(rec)
	            isPrinted = true
	        }
	        if tag != "" && tag == tagName {
	            logWriter.LogWrite(rec)
	            isPrinted = true
	        }
	    }
	}*/

	if !isPrinted {
		if lvl == INFO {
			_, _ = fmt.Fprintf(os.Stdout, formatLogRecord(defaultFormat, rec))
		} else if lvl > INFO {
			_, _ = fmt.Fprintf(os.Stderr, formatLogRecord(defaultFormat, rec))
		}
	}
}

func (p *loggerHandler) addLogFunc(lvl Level, withStack bool, tag string, logString string, src string) {
	rec := &logRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: logString,
	}

	p.print(rec, lvl, tag)
}

func (p *loggerHandler) addEmptyLine(lvl Level, tag string) {
	p.print(nil, lvl, tag)
}

// Load XML configuration
func (p *loggerHandler) loadConfiguration(filename string) {
	p.close()

	// Open the configuration file
	fd, err := os.Open(filename)
	if err != nil {
		printlnIO(os.Stderr, "ERROR", "open file:%s err: %s", filename, err.Error())
		os.Exit(1)
	}

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		printlnIO(os.Stderr, "ERROR", "read file all:%s err: %s", filename, err.Error())
		os.Exit(1)
	}

	xc := new(xmlLoggerConfig)
	if err := xml.Unmarshal(contents, xc); err != nil {
		printlnIO(os.Stderr, "ERROR", "xml.Unmarshal err: %s", err.Error())
		os.Exit(1)
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, xmlFilter := range xc.Filter {

		// Check required children
		if xmlFilter.Enabled == "" {
			printlnIO(os.Stderr, "ERROR", "log filter property:enabled not found")
			os.Exit(1)

		} else if xmlFilter.Enabled == "false" {
			continue
		}

		if xmlFilter.Tag == "" {
			printlnIO(os.Stderr, "ERROR", "log filter child:<tag> not found")
			os.Exit(1)

		} else if _, ok := p.logWriterMap[xmlFilter.Tag]; ok {
			printlnIO(os.Stderr, "ERROR", "log filter child:<tag>'s value repeat")
			os.Exit(1)
		}

		if xmlFilter.Type == "" {
			printlnIO(os.Stderr, "ERROR", "log filter child:<type> not found")
			os.Exit(1)
		}

		if xmlFilter.Level == "" {
			printlnIO(os.Stderr, "ERROR", "log filter child:<level> not found")
			os.Exit(1)
		}

		var lvl Level
		switch xmlFilter.Level {
		case "DEBUG":
			lvl = DEBUG
		case "INFO":
			lvl = INFO
		case "WARNING":
			lvl = WARNING
		case "ERROR":
			lvl = ERROR
		default:
			printlnIO(os.Stderr, "ERROR", "unsupported filter child:<level>'s value: %s", xmlFilter.Level)
			os.Exit(1)
		}

		logWriter, err := logWriter(nil), error(nil)
		switch xmlFilter.Type {
		case "console":
			logWriter, err = xmlToConsoleLogWriter(lvl, xmlFilter.Property)
		case "file":
			logWriter, err = xmlToFileLogWriter(xmlFilter.Tag, lvl, xmlFilter.Property)
		default:
			printlnIO(os.Stderr, "ERROR", "unsupported filter child:<type>'s value: %s", xmlFilter.Type)
			os.Exit(1)
		}

		if err != nil {
			printlnIO(os.Stderr, "ERROR", err.Error())
			os.Exit(1)

		} else if p.defaultLogFilePath == "" && xmlFilter.Type == "file" {
			fileLogWriter := logWriter.(*FileLogWriter)
			pathIndex := strings.LastIndex(fileLogWriter.GetFilename(), "/")
			p.defaultLogFilePath = fileLogWriter.GetFilename()[0:pathIndex]
		}

		p.logWriterMap[xmlFilter.Tag] = logWriter
	}
}

func (p *loggerHandler) addFileLoggerIfNotExist(tag string, lv Level, prop *LogProperty) (isExist bool) {

	p.lock.Lock()
	defer p.lock.Unlock()

	if _, ok := p.logWriterMap[tag]; !ok {
		if flw, err := NewFileLogWriter(tag, lv, prop.Filename, prop.Rotate, prop.KeepDay); err == nil {
			flw.SetFormat(prop.Format)
			flw.SetRotateLines(prop.MaxLines)
			flw.SetRotateSize(prop.Maxsize)
			flw.SetRotateDaily(prop.Daily)
			flw.SetPrivate(prop.Private)
			p.logWriterMap[tag] = flw
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func getSrcAndMsg(runtimeSkip int, withStack bool, format string, args ...interface{}) (string, string) {

	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(runtimeSkip)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	msg := bytes.Buffer{}
	if len(args) > 0 {
		msg.WriteString(fmt.Sprintf(format, args...))
	} else {
		msg.WriteString(format)
	}

	// 堆栈信息
	if withStack {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]

		msg.WriteByte('\n')
		if runtimeSkip <= 0 {
			msg.Write(buf) // 不跳过
		} else {
			// 第一行为goroutine标识，从第2行起，每2行是一个层级信息
			stack := strings.Split(string(buf), "\n")
			if stackLine := len(stack); stackLine > runtimeSkip*2+1 {
				newStack := make([]string, stackLine-runtimeSkip*2)
				newStack[0] = stack[0]
				copy(newStack[1:], stack[1+runtimeSkip*2:])
				msg.WriteString(strings.Join(newStack, "\n"))
			} else {
				msg.Write(buf) // should not happen
			}
		}
	}
	return src, msg.String()
}

func xmlToConsoleLogWriter(lvl Level, props []xmlProperty) (*ConsoleLogWriter, error) {

	format := "[%D %T] [%L] (%S) %M"

	for _, prop := range props {
		switch prop.Name {
		case "format":
			format = strings.Trim(prop.Value, " \r\n")
		default:
			return nil, fmt.Errorf("unsupported filter property: %s", prop.Name)
		}
	}

	console := NewConsoleLogWriter(lvl)
	console.SetFormat(format)
	return console, nil
}

// Parse a number with K/M/G suffixes based on thousands (1000) or 2^10 (1024)
func strToNumSuffix(str string, multiple int) int {
	num := 1
	if len(str) > 1 {
		switch str[len(str)-1] {
		case 'G', 'g':
			num *= multiple
			fallthrough
		case 'M', 'm':
			num *= multiple
			fallthrough
		case 'K', 'k':
			num *= multiple
			str = str[0 : len(str)-1]
		}
	}
	parsed, _ := strconv.Atoi(str)
	return parsed * num
}

func xmlToFileLogWriter(tag string, lvl Level, props []xmlProperty) (*FileLogWriter, error) {

	file := ""
	format := "[%D %T] [%L] (%S) %M"
	maxLines := 0
	maxSize := 0
	daily := false
	rotate := false
	private := false
	keepDay := int64(0)

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "filename":
			file = strings.Trim(prop.Value, " \r\n")
		case "format":
			format = strings.Trim(prop.Value, " \r\n")
		case "maxlines":
			maxLines = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000)
		case "maxsize":
			maxSize = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1024)
		case "daily":
			daily = strings.Trim(prop.Value, " \r\n") != "false"
		case "rotate":
			rotate = strings.Trim(prop.Value, " \r\n") != "false"
		case "private":
			private = strings.Trim(prop.Value, " \r\n") != "false"
		case "keepDay":
			keepDay = int64(strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000))
		default:
			return nil, fmt.Errorf("unsupported property: %s", prop.Name)
		}
	}

	// Check properties
	if len(file) == 0 {
		return nil, errors.New("missing property: filename")
	}

	if flw, err := NewFileLogWriter(tag, lvl, file, rotate, keepDay); err == nil {
		flw.SetFormat(format)
		flw.SetRotateLines(maxLines)
		flw.SetRotateSize(maxSize)
		flw.SetRotateDaily(daily)
		flw.SetPrivate(private)
		return flw, nil
	} else {
		return nil, err
	}

}

func printlnIO(ioWriter io.Writer, typ string, format string, args ...interface{}) {
	format = "[%s] [%s] [log4j] " + format
	args = append([]interface{}{time.Now().Format("2006/01/02 15:04:05"), typ}, args...)
	_, _ = fmt.Fprintln(ioWriter, fmt.Sprintf(format, args...))
}
