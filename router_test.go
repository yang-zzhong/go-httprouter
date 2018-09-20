package httprouter

import (
	"log"
	"net/http"
	"net/url"
	. "testing"
)

type RW struct {
	Code int
}

func (rw *RW) Header() http.Header {
	return http.Header{}
}

func (rw *RW) Write(msg []byte) (size int, err error) {
	msg = []byte{}
	size = 0
	return
}

func (rw *RW) WriteHeader(status int) {
	rw.Code = status
}

func getWriter() http.ResponseWriter {
	return &RW{200}
}

type middleware1 struct{}

func (mid *middleware1) Before(_ *ResponseWriter, _ *Request) bool {
	log.Print("middle 1 before")
	beforeMiddleware1Exec = true
	return true
}

func (mid *middleware1) After(_ *ResponseWriter, _ *Request) bool {
	log.Print("middle 1 after")
	afterMiddleware1Exec = true
	return true
}

type middleware2 struct{}

func (mid *middleware2) Before(_ *ResponseWriter, _ *Request) bool {
	log.Print("middle 2 before")
	beforeMiddleware2Exec = true
	return true
}

func (mid *middleware2) After(_ *ResponseWriter, _ *Request) bool {
	log.Print("middle 2 after")
	afterMiddleware2Exec = true
	return true
}

var router *Router
var helloWorldExec bool
var apiHelloWorldExec bool
var beforeMiddleware1Exec bool
var afterMiddleware1Exec bool
var beforeMiddleware2Exec bool
var afterMiddleware2Exec bool
var withMiddlewareExec bool
var params bool

func _beforeFile(_ *ResponseWriter, _ string) bool {
	log.Print("before file")
	return true
}

func init() {
	router = NewRouter()
	router.BeforeFile = _beforeFile

	helloWorldExec = false
	apiHelloWorldExec = false
	beforeMiddleware1Exec = false
	beforeMiddleware2Exec = false
	afterMiddleware1Exec = false
	afterMiddleware2Exec = false
	withMiddlewareExec = false
	params = false
	router.OnGet("/hello-world", func(w *ResponseWriter, req *Request) {
		helloWorldExec = true
	})
	router.Group("/api", []Middleware{}, func(router *Router) {
		router.Get("/hello-world", func(w *ResponseWriter, req *Request) {
			apiHelloWorldExec = true
		})
	})
	router.OnGet("/users/:name", func(w *ResponseWriter, req *Request) {
		if req.Bag.Get("name") == "young" {
			params = true
		}
	})
	router.Group("", []Middleware{&middleware1{}}, func(router *Router) {
		router.Group("", []Middleware{&middleware2{}}, func(router *Router) {
			router.OnGet("/middleware", func(w *ResponseWriter, req *Request) {
				withMiddlewareExec = true
			})
		})
	})
}

func getRequest(method string, path string) *http.Request {
	u := new(url.URL)
	u.Path = path
	req := new(http.Request)
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
	if !beforeMiddleware1Exec || !beforeMiddleware2Exec || !afterMiddleware1Exec || !afterMiddleware2Exec {
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
