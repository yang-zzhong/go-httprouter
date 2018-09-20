package httprouter

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// wrap http.ResponseWriter, provide some useful functions
type ResponseWriter struct {
	StatusCode int
	writer     http.ResponseWriter
	headers    map[string]string
	content    []byte
}

// new response writer
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{200, w, make(map[string]string), nil}
}

// get http.ResponseWriter
func (rw *ResponseWriter) UnderlyingWriter() http.ResponseWriter {
	return rw.writer
}

// set http status code
func (rw *ResponseWriter) WithStatusCode(statusCode int) *ResponseWriter {
	rw.StatusCode = statusCode
	return rw
}

// get response headers you setted
func (rw *ResponseWriter) Headers() map[string]string {
	return rw.headers
}

// get response body
func (rw *ResponseWriter) Body() []byte {
	return rw.content
}

// add http header to respond, if same key assign value many times, last time will be effective
func (rw *ResponseWriter) WithHeader(key, val string) *ResponseWriter {
	rw.headers[key] = val
	return rw
}

// read content to set as body from io.Reader
func (rw *ResponseWriter) Read(reader io.Reader) {
	reader.Read(rw.content)
}

// output interface{} as json body, legacy, abundon future, use WriteJson instead
func (rw *ResponseWriter) Json(content interface{}) {
	rw.WithHeader("Content-Type", "text/json")
	var err error
	if rw.content, err = json.Marshal(content); err != nil {
		panic(err)
	}
}

// output string as body, legacy, abundon future, use WriteString instead
func (rw *ResponseWriter) String(content string) {
	rw.WithHeader("Content-Type", "text/plain")
	rw.content = []byte(content)
}

// output []byte
func (rw *ResponseWriter) Write(content []byte) (int, error) {
	rw.content = content
	return len(rw.content), nil
}

// output interface{} as json body
func (rw *ResponseWriter) WriteJson(content interface{}) {
	rw.Json(content)
}

// output string as body
func (rw *ResponseWriter) WriteString(content string) {
	rw.String(content)
}

// output file as body
func (rw *ResponseWriter) WriteFile(pathfile string) {
	if content, err := ioutil.ReadFile(pathfile); err != nil {
		panic(err)
	} else if _, err := rw.Write(content); err != nil {
		panic(err)
	}
	if contentType := guessContentType(pathfile); contentType != "" {
		rw.WithHeader("Content-Type", contentType)
	} else {
		rw.WithHeader("Content-Type", "text/html")
	}
}

// output a server error
func (rw *ResponseWriter) InternalError(err error) {
	rw.WithStatusCode(500).String(err.Error())
}

// output result
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
