package httprouter

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strings"
)

// wrap http.ResponseWriter, provide some useful functions
type Response struct {
	statusCode int
	headers    map[string]string
	body       io.Reader
	writer     http.ResponseWriter
}

// new response writer
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{200, make(map[string]string), nil, w}
}

func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}

// set http status code
func (r *Response) WithStatus(statusCode int) *Response {
	r.statusCode = statusCode
	return r
}

// get response headers you setted
func (r *Response) Headers() map[string]string {
	return r.headers
}

// get response body
func (r *Response) Body() io.Reader {
	return r.body
}

// add http header to respond, if same key assign value many times, last time will be effective
func (r *Response) WithHeader(key, val string) *Response {
	r.headers[key] = val
	return r
}

// read content to set as body from io.Reader
func (r *Response) WithBody(body io.Reader) {
	r.body = body
}

func (r *Response) WithString(content string) {
	r.WithHeader("Content-Type", "text/plain")
	r.body = strings.NewReader(content)
}

func (r *Response) WithFile(p string) error {
	body, err := os.Open(p)
	if err != nil {
		return err
	}
	r.body = body
	ct := guessContentType(p)
	if ct != "" {
		r.WithHeader("Content-Type", ct)
	}

	return nil
}

// output a server error
func (r *Response) InternalError(err error) {
	r.WithStatus(500).WithString(err.Error())
}

// output result
func (r *Response) Flush(req *http.Request) error {
	w := r.writer
	for key, val := range r.headers {
		w.Header().Set(key, val)
	}
	if r.body == nil {
		w.WriteHeader(r.statusCode)
		return nil
	}
	ae := []byte(req.Header.Get("Accept-Encoding"))
	if bytes.Index(ae, []byte("gzip")) != -1 {
		w.WriteHeader(r.statusCode)
		io.Copy(w, r.body)
		return nil
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(r.statusCode)

	var buf bytes.Buffer
	z := gzip.NewWriter(&buf)
	defer z.Close()
	io.Copy(z, r.body)
	z.Flush()
	io.Copy(w, &buf)

	return nil
}
