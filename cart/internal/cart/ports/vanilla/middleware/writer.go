package middleware

import (
	"bytes"
	"net/http"
)

type statusWriter struct {
	w          http.ResponseWriter
	err        *bytes.Buffer
	statusCode int
}

func (sw *statusWriter) Header() http.Header {
	return sw.w.Header()
}

func (sw *statusWriter) Write(bytes []byte) (int, error) {
	sw.err.Write(bytes)
	return sw.w.Write(bytes)
}

func (sw *statusWriter) WriteHeader(statusCode int) {
	sw.statusCode = statusCode
	sw.w.WriteHeader(statusCode)
}
