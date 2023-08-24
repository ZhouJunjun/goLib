package log4j

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// This log writer sends output to a file
type FileLogWriter struct {
	level level
	tag   string

	logRecordCh chan *logRecord
	closeCh     chan bool

	// for del file
	timeTicker *time.Ticker

	// The opened file
	filename string
	file     *os.File

	// The logging format
	format string

	// Rotate at line count
	maxLines int
	curLines int

	// Rotate at size
	maxSize int
	curSize int

	// Rotate daily
	daily bool

	// create file ymd
	ymd string

	// Keep old logFiles (.001, .002, etc)
	rotate bool

	private bool

	keepDay int64
}

func NewFileLogWriter(tag string, level level, filename string, rotate bool, keepDay int64) (*FileLogWriter, error) {
	writer := &FileLogWriter{
		tag:         tag,
		level:       level,
		logRecordCh: make(chan *logRecord, LogBufferLength),
		closeCh:     make(chan bool),
		filename:    filename,
		format:      "[%D %T] [%L] (%S) %M",
		rotate:      rotate,
		keepDay:     keepDay,
	}

	// 设置了日志保存天数，定时删除过期日志
	if writer.keepDay > 0 {
		writer.timeTicker = time.NewTicker(time.Second * 60)
		go writer.delFile()
	}

	// try open log file
	if err := writer.openFile(); err != nil {
		return nil, fmt.Errorf("fileLogWriter[%s], openFile:%s fail, err:%s", tag, filename, err)
	}

	go writer.writeLog()

	printlnIO(os.Stdout, "INFO", "fileLogWriter[%s], create success, filename:%s", tag, filename)
	return writer, nil
}

// This is the FileLogWriter's output method
func (w *FileLogWriter) LogWrite(rec *logRecord) {
	w.logRecordCh <- rec
}

func (w *FileLogWriter) Close() {
	close(w.logRecordCh)
	printlnIO(os.Stdout, "INFO", "fileLogWriter[%s] closed log channel", w.tag)

	<-w.closeCh //等待 writeLog() 将日志全部写完后return
	printlnIO(os.Stdout, "INFO", "fileLogWriter[%s] is closed", w.tag)
}

func (w *FileLogWriter) writeLog() {

	defer func() {
		if w.file != nil {
			err := w.file.Close()
			printlnIO(os.Stdout, "INFO", "fileLogWrite[%s], close log file:%s, err:%+v", w.tag, w.filename, err)
			w.file = nil
		}
		w.closeCh <- true
	}()

	for {
		logRecord, isAlive := <-w.logRecordCh
		if !isAlive {
			printlnIO(os.Stdout, "INFO", "fileLogWriter[%s] log channel is empty", w.tag)
			return
		}

		if w.rotate && w.file != nil {
			w.tryMoveFile()
		}

		// 写入日志
		size, err := 0, error(nil)
		if w.file != nil {
			size, err = fmt.Fprint(w.file, formatLogRecord(w.format, logRecord))
		} else {
			// 程序启动后file是不为空的(执行openFile()失败,主程序会启动失败)；如果运行中file为空，可能是切割日志时关闭了file又无法重新打开
			size, err = fmt.Fprint(os.Stdout, formatLogRecord(w.format, logRecord))
		}

		if err != nil {
			printlnIO(os.Stderr, "ERROR", "fileLogWriter[%s] fmt.Fprint fail, err:%s", w.tag, err.Error())

		} else if w.rotate {
			if w.maxLines > 0 {
				w.curLines++
			}
			if w.maxSize > 0 {
				w.curSize += size
			}
		}
	}
}

func (w *FileLogWriter) openFile() error {

	pathIndex := strings.LastIndex(w.filename, "/")
	err := os.MkdirAll(string([]byte(w.filename)[0:pathIndex]), 0660)
	if err != nil {
		return err
	}

	// 以追加的方式打开文件
	file, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	w.file = file
	w.ymd = getYmd()
	w.curLines = 0
	w.curSize = 0
	return nil
}

func (w *FileLogWriter) closeFile() error {
	if err := w.file.Close(); err == nil {
		w.file = nil
		return nil
	} else {
		return err
	}
}

// 限制文件大小 或 行数 或 按日切割
func (w *FileLogWriter) tryMoveFile() {

	ymd := getYmd()
	if !((w.daily && ymd != w.ymd) || (w.maxLines > 0 && w.curLines >= w.maxLines) || (w.maxSize > 0 && w.curSize >= w.maxSize)) {
		return
	}

	// 检查文件是否存在
	if isExist := w.isFileExist(w.filename); !isExist {
		return //should not happen
	}

	for i := 0; i <= 999; i++ {

		// 逐个尝试新文件名，stdout.log.ymd or stdout.log.ymd.001 ~ stdout.log.ymd.999
		tmpFileName := w.filename + "." + w.ymd
		if i > 0 {
			tmpFileName += fmt.Sprintf("-%03d", i)
		}

		// 检查文件存在; 不存在, 把当前log改名字； stdout.log ===> stdout.log.ymd[.001]
		if isExist := w.isFileExist(tmpFileName); !isExist {

			if err := w.file.Close(); err != nil {
				// should not happen; 后续使用输出流写日志
				printlnIO(os.Stderr, "ERROR", "fileLogWriter[%s] close file:%s fail, err:%s", w.tag, w.file.Name(), err.Error())
				w.file = nil

			} else {
				w.file = nil

				if err := os.Rename(w.filename, tmpFileName); err != nil {
					printlnIO(os.Stderr, "ERROR", "fileLogWriter[%s] rename file fail, err:%s", w.tag, err.Error())
				}

				// 无论是否rename成功，再次打开文件(创建/追加)
				if err := w.openFile(); err != nil {
					printlnIO(os.Stderr, "ERROR", "fileLogWriter[%s] open file:%s fail, err:%s", w.tag, w.filename, err.Error())
				}
			}
			return
		}
	}
}

func (w *FileLogWriter) isFileExist(filePath string) bool {
	fileInfo, err := os.Lstat(filePath)
	return err == nil && fileInfo != nil
}

// 清除n天前的日志
func (w *FileLogWriter) delFile() {

	for {
		<-w.timeTicker.C

		// log文件夹位置
		pathIndex := strings.LastIndex(w.filename, "/")
		path := w.filename[0:pathIndex]

		if folder, err := ioutil.ReadDir(path); err != nil {
			printlnIO(os.Stderr, "ERROR", "fileLogWriter[%s] readDir:%s fail, err:%s", w.tag, path, err.Error())
			return

		} else {
			timeNow := time.Now().Unix()
			for _, file := range folder {
				if !file.IsDir() && file.ModTime().Unix()+86400*w.keepDay < timeNow {

					filePath := path + "/" + file.Name()
					if filePath == w.filename { //正在打的日志不删
						continue
					}

					if strings.HasPrefix(filePath, w.filename) {
						if err := os.Remove(filePath); err != nil {
							printlnIO(os.Stderr, "ERROR", "fileLogWriter[%s] remove:%s fail, err:%s", w.tag, filePath, err.Error())
						} else {
							printlnIO(os.Stdout, "INFO", "fileLogWriter[%s] remove:%s success", w.tag, filePath)
						}
					}
				}
			}
		}
	}
}

func (w *FileLogWriter) IsPrivate() bool {
	return w.private
}

func (w *FileLogWriter) GetLevel() level {
	return w.level
}

func (w *FileLogWriter) SetFormat(format string) *FileLogWriter {
	w.format = format
	return w
}

func (w *FileLogWriter) SetRotateLines(maxLines int) *FileLogWriter {
	w.maxLines = maxLines
	return w
}

func (w *FileLogWriter) SetRotateSize(maxsize int) *FileLogWriter {
	w.maxSize = maxsize
	return w
}

func (w *FileLogWriter) SetRotateDaily(daily bool) *FileLogWriter {
	w.daily = daily
	return w
}

func (w *FileLogWriter) SetRotate(rotate bool) *FileLogWriter {
	w.rotate = rotate
	return w
}

func (w *FileLogWriter) SetPrivate(private bool) *FileLogWriter {
	w.private = private
	return w
}

func (w *FileLogWriter) GetFilename() string {
	return w.filename
}

func getYmd() string {
	return time.Now().Format("20060102")
}
