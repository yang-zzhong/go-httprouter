package httprouter

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
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

func (rw *ResponseWriter) Flush(req *http.Request, w http.ResponseWriter) {
	zipEnabled := true
	for key, val := range rw.headers {
		w.Header().Set(key, val)
		if key == "Content-Type" && bytes.Index([]byte(key), []byte("image")) != 1 {
			zipEnabled = false
		}
	}
	acceptEncoding := []byte(req.Header.Get("Accept-Encoding"))
	w.WriteHeader(rw.statusCode)
	if zipEnabled && bytes.Index(acceptEncoding, []byte("gzip")) != -1 {
		w.Header().Set("Content-Encoding", "gzip")
		z := gzip.NewWriter(w)
		defer z.Close()
		if _, err := z.Write(rw.content); err != nil {
			log.Print(err)
		}
		z.Flush()
	} else {
		w.Write(rw.content)
	}
}
