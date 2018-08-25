package httprouter

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
)

type ResponseWriter struct {
	StatusCode int
	writer     http.ResponseWriter
	headers    map[string]string
	content    []byte
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{200, w, make(map[string]string), nil}
}

func (rw *ResponseWriter) UnderlyingWriter() http.ResponseWriter {
	return rw.writer
}

func (rw *ResponseWriter) WithStatusCode(statusCode int) *ResponseWriter {
	rw.StatusCode = statusCode
	return rw
}

func (rw *ResponseWriter) Headers() map[string]string {
	return rw.headers
}

func (rw *ResponseWriter) Body() []byte {
	return rw.content
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

func (rw *ResponseWriter) Flush(req *http.Request) error {
	w := rw.writer
	for key, val := range rw.headers {
		w.Header().Set(key, val)
	}
	ae := []byte(req.Header.Get("Accept-Encoding"))
	if bytes.Index(ae, []byte("gzip")) == -1 {
		w.WriteHeader(rw.StatusCode)
		w.Write(rw.content)
		return nil
	}
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(rw.StatusCode)
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
