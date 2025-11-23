package webUtil

import (
	"compress/gzip"
	"net/http"
)

func NewGZipWriter(writer http.ResponseWriter) http.ResponseWriter {
	writer.Header().Set("Content-Encoding", "gzip")
	return &gzipWriter{
		originWriter: writer,
		gzipWriter:   gzip.NewWriter(writer),
	}
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
	return p.gzipWriter.Close()
}
