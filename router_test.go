package httprouter

import (
	helper "github.com/yang-zzhong/go-helpers"
	. "net/http"
	"net/url"
	. "testing"
)

type RW struct {
	Code int
}

func (rw *RW) Header() Header {
	return Header{}
}

func (rw *RW) Write(msg []byte) (size int, err error) {
	msg = []byte{}
	size = 0
	return
}

func (rw *RW) WriteHeader(status int) {
	rw.Code = status
}

func getWriter() ResponseWriter {
	return &RW{200}
}

var router *Router
var helloWorldExec bool
var apiHelloWorldExec bool
var middlewareExec bool
var withMiddlewareExec bool
var params bool

func init() {
	router = CreateRouter(".", []string{"index.html"})
	helloWorldExec = false
	apiHelloWorldExec = false
	middlewareExec = false
	withMiddlewareExec = false
	params = false
	router.Get("/hello-world", func(w ResponseWriter, req *Request, _ *helper.P) {
		helloWorldExec = true
	})
	router.Group("/api", NewMs(), func(router *Router) {
		router.Get("/hello-world", func(w ResponseWriter, req *Request, _ *helper.P) {
			apiHelloWorldExec = true
		})
	})
	router.Get("/users/:name", func(w ResponseWriter, req *Request, p *helper.P) {
		if p.Get("name") == "young" {
			params = true
		}
	})
	middle1 := func(w ResponseWriter, req *Request) bool {
		middlewareExec = true

		return true
	}
	ms := NewMs()
	ms.Append(middle1)
	router.Group("", ms, func(router *Router) {
		router.Get("/middleware", func(w ResponseWriter, req *Request, _ *helper.P) {
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
	writer := getWriter()
	router.ServeHTTP(writer, getRequest("GET", "/hello-world"))
	if !helloWorldExec {
		t.Error("hello-world fail")
	}
	router.ServeHTTP(writer, getRequest("GET", "/api/hello-world"))
	if !apiHelloWorldExec {
		t.Error("api/hello-world fail")
	}
	router.ServeHTTP(writer, getRequest("GET", "/middleware"))
	if !middlewareExec {
		t.Error("middleware fail")
	}
	if !withMiddlewareExec {
		t.Error("after middleware fail")
	}
	router.ServeHTTP(writer, getRequest("GET", "/users/young"))
	if !params {
		t.Error("params fail")
	}
	router.ServeHTTP(writer, getRequest("GET", "/not-found"))
	if writer.(*RW).Code != 404 {
		t.Error("not found fail")
	}

	router.ServeHTTP(writer, getRequest("POST", "/hello-world"))
	if writer.(*RW).Code != 405 {
		t.Error("not allowed fail")
	}
	router.ServeHTTP(writer, getRequest("GET", "/README.md"))
	if writer.(*RW).Code == 404 {
		t.Error("default router fail")
	}
}
