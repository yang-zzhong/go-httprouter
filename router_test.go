package httprouter

import (
	. "net/http"
	"net/url"
	. "testing"
)

type RW struct{}

func (rw *RW) Header() Header {
	return Header{}
}

func (rw *RW) Write(msg []byte) (size int, err error) {
	msg = []byte{}
	size = 0
	return
}

func (rw *RW) WriteHeader(status int) {
	status = 200
}

func getWriter() ResponseWriter {
	return new(RW)
}

var router *Router
var helloWorldExec bool
var apiHelloWorldExec bool
var middlewareExec bool
var withMiddlewareExec bool
var params bool

func init() {
	router = New("/tmp", []string{"index.html"})
	helloWorldExec = false
	apiHelloWorldExec = false
	middlewareExec = false
	withMiddlewareExec = false
	params = false
	router.Get("/hello-world", func(w ResponseWriter, req *Request, _ map[string]string) {
		helloWorldExec = true
	})
	router.Group("/api", []Middleware{}, func(router *Router) {
		router.Get("/hello-world", func(w ResponseWriter, req *Request, _ map[string]string) {
			apiHelloWorldExec = true
		})
	})
	router.Get("/users/:name", func(w ResponseWriter, req *Request, p map[string]string) {
		if p["name"] == "young" {
			params = true
		}
	})
	middle1 := func(w ResponseWriter, req *Request) bool {
		middlewareExec = true

		return true
	}
	router.Group("", []Middleware{middle1}, func(router *Router) {
		router.Get("/middleware", func(w ResponseWriter, req *Request, _ map[string]string) {
			withMiddlewareExec = true
		})
	})
}

func getRequest(method string, path string) *Request {
	u := new(url.URL)
	u.Path = path
	req := new(Request)
	req.Method = method
	req.URL = u

	return req
}

func TestRoute(t *T) {
	router.ServeHTTP(getWriter(), getRequest("GET", "/hello-world"))
	if !helloWorldExec {
		t.Error("hello-world fail")
	}
	router.ServeHTTP(getWriter(), getRequest("GET", "/api/hello-world"))
	if !apiHelloWorldExec {
		t.Error("api/hello-world fail")
	}
	router.ServeHTTP(getWriter(), getRequest("GET", "/middleware"))
	if !middlewareExec {
		t.Error("middleware fail")
	}
	if !withMiddlewareExec {
		t.Error("after middleware fail")
	}
	router.ServeHTTP(getWriter(), getRequest("GET", "/users/young"))
	if !params {
		t.Error("params fail")
	}
	router.ServeHTTP(getWriter(), getRequest("GET", "/not-found"))
}
