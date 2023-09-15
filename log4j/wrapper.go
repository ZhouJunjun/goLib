package log4j

import (
	"errors"
	"fmt"
)

var (
	globalHandler      *loggerHandler
	defaultRuntimeSkip = 3
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

func AddFileLoggerIfNotExist(tag string, lv level, logProperty *LogProperty) bool {
	return globalHandler.addFileLoggerIfNotExist(tag, lv, logProperty)
}

func GetLogFilePath() string {
	return globalHandler.getLogFilePath()
}

type LogBuffer interface {
	IsError() bool
	String() string
	PrintStack() bool
	RuntimeSkip(defaultSkip int) int
	Tag() string
}

func Log(logBuffer LogBuffer) {
	logTxt, skip := logBuffer.String(), logBuffer.RuntimeSkip(defaultRuntimeSkip)
	if logBuffer.IsError() {
		globalHandler.addLogString(skip, ERROR, logBuffer.PrintStack(), logBuffer.Tag(), logTxt)
	} else {
		globalHandler.addLogString(skip, INFO, false, logBuffer.Tag(), logTxt)
	}
}

func Debug(arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, DEBUG, false, "", first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(DEBUG, false, "", logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(defaultRuntimeSkip, DEBUG, false, "", "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, DEBUG, false, "", "%+v", arg0)
		}
	}
}

func DebugTag(tag string, arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, DEBUG, false, tag, first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(DEBUG, false, tag, logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(defaultRuntimeSkip, DEBUG, false, tag, "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, DEBUG, false, tag, "%+v", arg0)
		}
	}
}

func Info(arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, INFO, false, "", first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(INFO, false, "", logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(defaultRuntimeSkip, INFO, false, "", "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, INFO, false, "", "%+v", arg0)
		}
	}
}

func InfoTag(tag string, arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, INFO, false, tag, first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(INFO, false, tag, logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(defaultRuntimeSkip, INFO, false, tag, "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, INFO, false, tag, "%+v", arg0)
		}
	}
}

func Warn(arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, WARNING, false, "", first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(WARNING, false, "", logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(defaultRuntimeSkip, WARNING, false, "", "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, WARNING, false, "", "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func WarnTag(tag string, arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, WARNING, false, tag, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(WARNING, false, tag, logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(defaultRuntimeSkip, WARNING, false, tag, "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, WARNING, false, tag, "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func Error(arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, ERROR, false, "", first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(ERROR, false, "", logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(defaultRuntimeSkip, ERROR, false, "", "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, ERROR, false, "", "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func ErrorTag(tag string, arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, ERROR, false, tag, first, args...)
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(ERROR, false, tag, logString, src)
	default:
		if args != nil {
			globalHandler.addLogString(defaultRuntimeSkip, ERROR, false, tag, "%+v", append([]interface{}{arg0}, args))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, ERROR, false, tag, "%+v", arg0)
		}
	}
}

// error log with stack info
func ErrorStack(arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, ERROR, true, "", first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(ERROR, true, "", logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(defaultRuntimeSkip, ERROR, true, "", "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, ERROR, true, "", "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func ErrorTagStack(tag string, arg0 interface{}, args ...interface{}) error {
	switch first := arg0.(type) {
	case string:
		globalHandler.addLogString(defaultRuntimeSkip, ERROR, true, tag, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() (log string, src string):
		logString, src := first()
		globalHandler.addLogFunc(ERROR, true, tag, logString, src)
		return errors.New(logString)
	default:
		if args != nil {
			slice := append([]interface{}{arg0}, args)
			globalHandler.addLogString(defaultRuntimeSkip, ERROR, true, tag, "%+v", slice)
			return errors.New(fmt.Sprintf("%+v", slice))
		} else {
			globalHandler.addLogString(defaultRuntimeSkip, ERROR, true, tag, "%+v", arg0)
			return errors.New(fmt.Sprintf("%+v", arg0))
		}
	}
}

func EmptyLine(lvl level, tag string) {
	globalHandler.addEmptyLine(lvl, tag)
}
