package httprouter

import (
	helper "github.com/yang-zzhong/go-helpers"
	"io"
	"io/ioutil"
	. "net/http"
	"os"
)

type HttpHandler func(ResponseWriter, *Request, *helper.P)
type GroupCall func(router *Router)

type Router struct {
	DocRoot string
	Indexes []string
	On404   HttpHandler
	configs []config
	ms      *Middlewares
	prefix  string
}

type config struct {
	method string
	path   string
	ms     *Middlewares
	call   HttpHandler
}

func onNotFound(w ResponseWriter, req *Request, _ *helper.P) {
	w.WriteHeader(StatusNotFound)
	io.WriteString(w, "not found")
}

func CreateRouter(docRoot string, indexes []string) *Router {
	router := new(Router)

	router.DocRoot = docRoot
	router.Indexes = indexes
	router.configs = []config{}
	router.ms = NewMs()
	router.prefix = ""
	router.On404 = onNotFound
	return router
}

func (router *Router) ServeHTTP(w ResponseWriter, req *Request) {
	for _, conf := range router.configs {
		matched, params := router.Match(conf.method, conf.path, req)
		if !matched {
			continue
		}
		if req.Method != conf.method {
			w.WriteHeader(StatusMethodNotAllowed)
			return
		}
		if conf.ms.Exec(w, req) {
			conf.call(w, req, params)
		}
		return
	}
	if req.Method == MethodGet {
		router.defaultRoute(w, req)
		return
	}
	w.WriteHeader(StatusNotFound)
}

func (router *Router) defaultRoute(w ResponseWriter, req *Request) {
	data, rerr := router.ReadFile(req.URL.Path)
	if rerr == nil {
		io.WriteString(w, (string)(data))
		return
	}
	pathfile, ierr := router.IndexFile(req.URL.Path)
	if ierr != nil {
		router.On404(w, req, helper.NewP())
		return
	}
	data, err := router.ReadFile(pathfile)
	if err != nil {
		router.On404(w, req, helper.NewP())
		return
	}

	io.WriteString(w, (string)(data))
}

func (router *Router) ReadFile(path string) (data []byte, err error) {
	pathfile := router.DocRoot + "/" + path
	if _, err = os.Stat(pathfile); err != nil {
		data = []byte{}
		return
	}

	return ioutil.ReadFile(pathfile)
}

func (router *Router) IndexFile(path string) (pathfile string, err error) {
	for _, index := range router.Indexes {
		pf := router.DocRoot + path + "/" + index
		if _, err = os.Stat(pf); err == nil {
			pathfile = pf
			return
		}
	}
	err = NewHE(StatusNotFound, "文件没有找到")
	return
}

func (router *Router) Match(method string, path string, req *Request) (m bool, p *helper.P) {
	m, p = newPath(path).match(req.URL.Path)
	return
}

func (router *Router) Get(path string, h HttpHandler) {
	router.Handle(MethodGet, path, h)
}

func (router *Router) Post(path string, h HttpHandler) {
	router.Handle(MethodPost, path, h)
}

func (router *Router) Put(path string, h HttpHandler) {
	router.Handle(MethodPut, path, h)
}

func (router *Router) Delete(path string, h HttpHandler) {
	router.Handle(MethodDelete, path, h)
}

func (router *Router) Patch(path string, h HttpHandler) {
	router.Handle(MethodPatch, path, h)
}

func (router *Router) Options(path string, h HttpHandler) {
	router.Handle(MethodOptions, path, h)
}

func (router *Router) Trace(path string, h HttpHandler) {
	router.Handle(MethodTrace, path, h)
}

func (router *Router) Connect(path string, h HttpHandler) {
	router.Handle(MethodConnect, path, h)
}

func (router *Router) Handle(method string, path string, h HttpHandler) {
	router.configs = append(
		router.configs, config{method, router.prefix + path, router.ms, h},
	)
}

func (router *Router) Group(prefix string, ms *Middlewares, grp GroupCall) {
	router.ms.Merge(ms)
	router.prefix += prefix
	grp(router)
	router.ms = NewMs()
	router.prefix = ""
}
