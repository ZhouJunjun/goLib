// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4j

import (
	"errors"
	"fmt"
)

const (
	// 跳过调用层级；不然会把本包内的调用路径也打出来
	RUNTIME_SKIP = 3
)

var (
	globalHandler *loggerHandler
)

func init() {
	globalHandler = newDefaultLogger(DEBUG)
}

func LoadConfiguration(filename string) {
	globalHandler.loadConfiguration(filename)
}

func Close() {
	globalHandler.close()
}

func CloseByTag(logTag string) {
	globalHandler.closeByTag(logTag)
}

func AddFileLoggerIfNotExist(tag string, lv Level, logProperty *LogProperty) bool {
	return globalHandler.addFileLoggerIfNotExist(tag, lv, logProperty)
}

func GetLogFilePath() string {
	return globalHandler.getLogFilePath()
}

type LogBuffer interface {
	String() string
	PrintStack() bool
	RuntimeSkip(defaultSkip int) int
	GetLogLevel() Level
}

func Log(logBuffer LogBuffer) {
	logTxt, skip := logBuffer.String(), logBuffer.RuntimeSkip(RUNTIME_SKIP)
	switch lv := logBuffer.GetLogLevel(); lv {
	case DEBUG, INFO, WARNING:
		globalHandler.addLogString(skip, lv, false, "", logTxt)
	default:
		globalHandler.addLogString(skip, lv, logBuffer.PrintStack(), "", logTxt)
	}
}

func LogIfError(logBuffer LogBuffer) {
	if logBuffer.GetLogLevel() == ERROR {
		logTxt, skip := logBuffer.String(), logBuffer.RuntimeSkip(RUNTIME_SKIP)
		globalHandler.addLogString(skip, ERROR, logBuffer.PrintStack(), "", logTxt)
	}
}

func LogTag(tag string, logBuffer LogBuffer) {
	logTxt, skip := logBuffer.String(), logBuffer.RuntimeSkip(RUNTIME_SKIP)
	switch lv := logBuffer.GetLogLevel(); lv {
	case DEBUG, INFO, WARNING:
		globalHandler.addLogString(skip, lv, false, tag, logTxt)
	default:
		globalHandler.addLogString(skip, lv, logBuffer.PrintStack(), tag, logTxt)
	}
}

func LogTagIfError(tag string, logBuffer LogBuffer) {
	if logBuffer.GetLogLevel() == ERROR {
		logTxt, skip := logBuffer.String(), logBuffer.RuntimeSkip(RUNTIME_SKIP)
		globalHandler.addLogString(skip, ERROR, logBuffer.PrintStack(), tag, logTxt)
	}
}

func Debug(arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, DEBUG, false, "", first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(DEBUG, false, "", logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(RUNTIME_SKIP, DEBUG, false, "", "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, DEBUG, false, "", "%+v", arg0)
		}
	}
}

func DebugTag(tag string, arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, DEBUG, false, tag, first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(DEBUG, false, tag, logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(RUNTIME_SKIP, DEBUG, false, tag, "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, DEBUG, false, tag, "%+v", arg0)
		}
	}
}

func Info(arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, INFO, false, "", first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(INFO, false, "", logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(RUNTIME_SKIP, INFO, false, "", "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, INFO, false, "", "%+v", arg0)
		}
	}
}

func InfoTag(tag string, arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, INFO, false, tag, first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(INFO, false, tag, logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(RUNTIME_SKIP, INFO, false, tag, "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, INFO, false, tag, "%+v", arg0)
		}
	}
}

func Warn(arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, WARNING, false, "", first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(WARNING, false, "", logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(RUNTIME_SKIP, WARNING, false, "", "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, WARNING, false, "", "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func WarnTag(tag string, arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, WARNING, false, tag, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(WARNING, false, tag, logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(RUNTIME_SKIP, WARNING, false, tag, "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, WARNING, false, tag, "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func Error(arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, ERROR, false, "", first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(ERROR, false, "", logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(RUNTIME_SKIP, ERROR, false, "", "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, ERROR, false, "", "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func ErrorTag(tag string, arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, ERROR, false, tag, first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(ERROR, false, tag, logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(RUNTIME_SKIP, ERROR, false, tag, "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, ERROR, false, tag, "%+v", arg0)
		}
	}
}

// error log with stack info
func ErrorStack(arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, ERROR, true, "", first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(ERROR, true, "", logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(RUNTIME_SKIP, ERROR, true, "", "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, ERROR, true, "", "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func ErrorTagStack(tag string, arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(RUNTIME_SKIP, ERROR, true, tag, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(ERROR, true, tag, logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(RUNTIME_SKIP, ERROR, true, tag, "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(RUNTIME_SKIP, ERROR, true, tag, "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func EmptyLine(lvl Level, tag string) {
	globalHandler.addEmptyLine(lvl, tag)
}
