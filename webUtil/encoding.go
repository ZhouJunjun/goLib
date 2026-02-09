package webUtil

import (
	"compress/gzip"
	"github.com/ZhouJunjun/goLib/log4j"
	"io"
	"net/http"
	"sync"
	"time"
)

var writerNum = 0
var GzipWriterPool = sync.Pool{
	New: func() interface{} {
		writerNum++
		return gzip.NewWriter(io.Discard) // 不用io.Discard，使用的时候reset
	},
}

func init() {
	go func() {
		log4j.Info("gzip writer num: %d", writerNum)
		time.Sleep(time.Second * 10)
	}()
}

func NewGZipWriter(writer http.ResponseWriter) *gzipWriter {
	writer.Header().Set("Content-Encoding", "gzip")
	// log4j.Info("new gzip writer")
	rw := &gzipWriter{
		originWriter: writer,
	}

	rw.gzipWriter = GzipWriterPool.Get().(*gzip.Writer)
	rw.gzipWriter.Reset(writer)
	return rw
}

type gzipWriter struct {
	originWriter http.ResponseWriter
	gzipWriter   *gzip.Writer
}

func (p *gzipWriter) Header() http.Header {
	return p.originWriter.Header()
}

func (p *gzipWriter) Write(data []byte) (int, error) {
	return p.gzipWriter.Write(data)
}

func (p *gzipWriter) WriteHeader(statusCode int) {
	p.originWriter.WriteHeader(statusCode)
}

func (p *gzipWriter) Close() error {
	// log4j.Info("close gzip writer")
	if err := p.gzipWriter.Close(); err != nil { // close，确保缓冲区flush
		return err
	} else {
		GzipWriterPool.Put(p.gzipWriter) // 放回池中
		return nil
	}
}
