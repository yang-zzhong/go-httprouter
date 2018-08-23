package httprouter

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
)

type ResponseWriter struct {
	statusCode int
	headers    map[string]string
	content    []byte
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{200, make(map[string]string), nil}
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
	rw.WithHeader("Content-Type", "text/json")
	rw.content, _ = json.Marshal(content)
}

func (rw *ResponseWriter) String(content string) {
	rw.WithHeader("Content-Type", "text/plain")
	rw.content = []byte(content)
}

func (rw *ResponseWriter) Write(content []byte) (int, error) {
	rw.content = content
	return len(rw.content), nil
}

func (rw *ResponseWriter) InternalError(err error) {
	rw.WithStatusCode(500).String(err.Error())
}

func (rw *ResponseWriter) Flush(req *http.Request, w http.ResponseWriter) error {
	for key, val := range rw.headers {
		w.Header().Set(key, val)
	}
	ae := []byte(req.Header.Get("Accept-Encoding"))
	if bytes.Index(ae, []byte("gzip")) == -1 {
		w.WriteHeader(rw.statusCode)
		w.Write(rw.content)
		return nil
	}
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(rw.statusCode)
	var buf bytes.Buffer
	z := gzip.NewWriter(&buf)
	defer z.Close()
	if _, err := z.Write(rw.content); err != nil {
		return err
	}
	z.Flush()
	w.Write(buf.Bytes())
	return nil
}
