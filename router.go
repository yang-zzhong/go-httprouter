package httprouter

import (
	helper "github.com/yang-zzhong/go-helpers"
	"io"
	. "net/http"
	"strconv"
)

const (
	Api       = "api"
	PathFile  = "pathfile"
	EntryFile = "entryfile"
)

type HttpHandler func(ResponseWriter, *Request, *helper.P)
type GroupCall func(router *Router)

type Router struct {
	Tries     []string
	DocRoot   string
	EntryFile string
	configs   []config
	ms        *Middlewares
	prefix    string
}

type config struct {
	method string
	path   string
	ms     *Middlewares
	call   HttpHandler
}

func NewRouter() *Router {
	router := new(Router)
	router.Tries = []string{Api, PathFile, EntryFile}
	router.DocRoot = "."
	router.EntryFile = "index.html"
	router.configs = []config{}
	router.ms = NewMs()
	router.prefix = ""

	return router
}

func (router *Router) ServeHTTP(w ResponseWriter, req *Request) {
	if req.Method == MethodGet {
		router.try(w, req)
		return
	}
	if router.tryApi(w, req) {
		return
	}
	w.WriteHeader(StatusNotFound)
}

func (router *Router) try(w ResponseWriter, req *Request) {
	for _, try := range router.Tries {
		switch try {
		case "api":
			if router.tryApi(w, req) {
				return
			}
		case "entryfile":
			if router.tryEntryFile(w, req) {
				return
			}
		case "pathfile":
			if router.tryPathFile(w, req) {
				return
			}
		}
	}
	w.WriteHeader(StatusNotFound)
}

func (router *Router) tryApi(w ResponseWriter, req *Request) bool {
	for _, conf := range router.configs {
		matched, params := router.Match(conf.method, conf.path, req)
		if !matched {
			continue
		}
		if req.Method != conf.method {
			w.WriteHeader(StatusMethodNotAllowed)
			return true
		}
		if conf.ms.Exec(w, req) {
			conf.call(w, req, params)
		}
		return true
	}

	return false
}

func (router *Router) tryEntryFile(w ResponseWriter, req *Request) bool {
	return router.tryFile(w, router.EntryFile)
}

func (router *Router) tryPathFile(w ResponseWriter, req *Request) bool {
	return router.tryFile(w, req.URL.Path)
}

func (router *Router) tryFile(w ResponseWriter, file string) bool {
	fh := newFileHandler(router.DocRoot)
	available, _ := fh.Available(file)
	if !available {
		return false
	}
	content, cerr := fh.Contents(file)
	if cerr != nil {
		return false
	}
	w.Header().Set("Content-Type", fh.ContentType(file))
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(StatusOK)

	io.WriteString(w, (string)(content))
	return true
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
