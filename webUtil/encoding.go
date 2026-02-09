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
	go func(num *int) {
		for {
			log4j.Info("gzip writer num: %d", *num)
			time.Sleep(time.Second * 10)
		}
	}(&writerNum)
}

func NewGZipWriter(writer http.ResponseWriter) *gzipWriter {
	writer.Header().Set("Content-Encoding", "gzip")
	// log4j.Info("new gzip writer")
	rw := &gzipWriter{
		originWriter: writer,
	}

	rw.gzipWriter = GzipWriterPool.Get().(*gzip.Writer)
	rw.gzipWriter.Reset(writer)
	// log4j.Info("[%p] http get gzip writer", rw.gzipWriter)
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
		_ = log4j.Error("gzip writer close error: %s", err.Error())
		return err
	} else {
		// log4j.Info("[%p] http return gzip writer", p.gzipWriter)
		GzipWriterPool.Put(p.gzipWriter) // 放回池中
		return nil
	}
}
