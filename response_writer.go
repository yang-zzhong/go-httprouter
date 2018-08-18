package httprouter

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
)

type ResponseWriter struct {
	w           http.ResponseWriter
	statusCode  int
	headers     map[string]string
	gzipEnabled bool
	content     []byte
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, 200, make(map[string]string), true, nil}
}

func (rw *ResponseWriter) WithStatusCode(statusCode int) *ResponseWriter {
	rw.statusCode = statusCode
	return rw
}

func (rw *ResponseWriter) WithHeader(key, val string) *ResponseWriter {
	rw.headers[key] = val
	return rw
}

func (rw *ResponseWriter) Read(reader io.Reader) {
	reader.Read(rw.content)
}

func (rw *ResponseWriter) Json(content interface{}) {
	rw.content, _ = json.Marshal(content)
}

func (rw *ResponseWriter) String(content string) {
	rw.content = []byte(content)
}

func (rw *ResponseWriter) Write(content []byte) (int, error) {
	rw.content = content
	return len(rw.content), nil
}

func (rw *ResponseWriter) InternalError(err error) {
	rw.WithStatusCode(500).String(err.Error())
}

func (rw *ResponseWriter) Flush(req *http.Request) {
	rw.w.WriteHeader(rw.statusCode)
	for key, val := range rw.headers {
		rw.w.Header().Set(key, val)
	}
	if rw.gzipEnabled {
		rw.w.Header().Set("Content-Encoding", "gzip")
		w := gzip.NewWriter(rw.w)
		w.Write(rw.content)
	} else {
		rw.w.Write(rw.content)
	}
}
