package compresser

import (
	"github.com/gin-gonic/gin"
	"io"
)

// Compressor middleware compressor
type Compressor struct {
	gin.ResponseWriter
	Writer io.Writer
}

// NewCompressor return new middleware compressor
func NewCompressor() gin.HandlerFunc {
	return newCompressHandler().Handle
}

// Write to writer
func (c *Compressor) Write(b []byte) (int, error) {
	return c.Writer.Write(b)
}

// WriteString to writer
func (c *Compressor) WriteString(s string) (int, error) {
	c.Header().Del("Content-Length")
	return c.Writer.Write([]byte(s))
}

// WriteHeader write to header
func (c *Compressor) WriteHeader(code int) {
	c.Header().Del("Content-Length")
	c.ResponseWriter.WriteHeader(code)
}
