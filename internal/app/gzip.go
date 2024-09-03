package app

import (
	"compress/gzip"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	contentType := c.w.Header().Get("Content-Type")
	if strings.Contains(contentType, echo.MIMEApplicationJSON) || strings.Contains(contentType, echo.MIMETextHTML) {
		return c.zw.Write(p)
	}
	return c.w.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	contentType := c.w.Header().Get("Content-Type")
	if statusCode < 500 && (strings.Contains(contentType, echo.MIMEApplicationJSON) || strings.Contains(contentType, echo.MIMETextHTML)) {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	contentType := c.w.Header().Get("Content-Type")
	if strings.Contains(contentType, echo.MIMEApplicationJSON) || strings.Contains(contentType, echo.MIMETextHTML) {
		return c.zw.Close()
	}
	return nil
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func NewGZIPMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resp := c.Response()
			w := resp.Writer

			req := c.Request()
			acceptEncoding := req.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip") {
				cw := &compressWriter{w: w, zw: gzip.NewWriter(w)}
				defer cw.Close()
				resp.Writer = cw
			}

			contentEncoding := req.Header.Get("Content-Encoding")
			if strings.Contains(contentEncoding, "gzip") {
				cr, err := newCompressReader(req.Body)
				if err != nil {
					return c.NoContent(http.StatusBadRequest)
				}
				defer cr.Close()
				req.Body = cr
			}
			return next(c)
		}
	}
}
