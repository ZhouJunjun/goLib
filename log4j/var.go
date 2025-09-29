package log4j

const (
    // LogBufferLength specifies how many log messages a particular log4go
    // logger can buffer at a time before writing them.
    LogBufferLength = 32

    defaultFormat = "[%D %T] [%L] (%S) %M"
)

type Level int

const (
    DEBUG Level = iota
    INFO
    WARNING
    ERROR
)

var (
    levelStrings = [...]string{"DEBG", "INFO", "WARN", "EROR"}
)

func (l Level) String() string {
    if l < 0 || int(l) > len(levelStrings) {
        return "UNKNOWN"
    }
    return levelStrings[int(l)]
}
