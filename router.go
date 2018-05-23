package httprouter

import (
	helper "github.com/yang-zzhong/go-helpers"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	Api       = "api"
	PathFile  = "pathfile"
	EntryFile = "entryfile"
)

type HttpHandler func(http.ResponseWriter, *Request, *helper.P)
type GroupCall func(router *Router)
type ResponseHeaderHandler func(http.ResponseWriter)

type Router struct {
	Tries                []string
	DocRoot              string
	EntryFile            string
	On404                HttpHandler
	configs              []config
	Logger               *log.Logger
	HandleResponseHeader ResponseHeaderHandler
	ms                   *Middlewares
	prefix               string
}

type config struct {
	method string
	path   string
	ms     *Middlewares
	call   HttpHandler
}

func onNotFound(w http.ResponseWriter, req *Request, _ *helper.P) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "not found")
}

func NewRouter() *Router {
	router := new(Router)
	router.Tries = []string{Api, PathFile, EntryFile}
	router.DocRoot = "."
	router.EntryFile = "index.html"
	router.configs = []config{}
	router.ms = NewMs()
	router.prefix = ""
	router.Logger = log.New(os.Stdout, "Http Router -> ", log.Lshortfile)
	router.On404 = onNotFound
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.Logger.Printf("%v\t%s\t%v\t%v", req.Proto, req.URL.Path, req.Header, req.Body)
	if req.Method == http.MethodGet {
		router.try(w, req)
		return
	}
	if router.tryApi(w, req) {
		return
	}
	router.On404(w, &Request{req}, helper.NewP())
}

func (router *Router) try(w http.ResponseWriter, req *http.Request) {
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
	router.On404(w, &Request{req}, helper.NewP())
}

func (router *Router) tryApi(w http.ResponseWriter, req *http.Request) bool {
	methodNotAllowed := false
	for _, conf := range router.configs {
		matched, params := router.Match(conf.method, conf.path, req)
		if !matched {
			continue
		}
		if req.Method != conf.method {
			methodNotAllowed = true
			continue
		}
		router.HandleResponseHeader(w)
		if conf.ms.Exec(w, &Request{req}, params) {
			conf.call(w, &Request{req}, params)
		}
		return true
	}
	if methodNotAllowed {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return true
	}

	return false
}

func (router *Router) tryEntryFile(w http.ResponseWriter, req *http.Request) bool {
	return router.tryFile(w, router.EntryFile)
}

func (router *Router) tryPathFile(w http.ResponseWriter, req *http.Request) bool {
	return router.tryFile(w, req.URL.Path)
}

func (router *Router) tryFile(w http.ResponseWriter, file string) bool {
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
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, (string)(content))
	return true
}

func (router *Router) Match(method string, path string, req *http.Request) (m bool, p *helper.P) {
	m, p = newPath(path).match(req.URL.Path)
	return
}

func (router *Router) Get(path string, h HttpHandler) {
	router.Handle(http.MethodGet, path, h)
}

func (router *Router) Post(path string, h HttpHandler) {
	router.Handle(http.MethodPost, path, h)
}

func (router *Router) Put(path string, h HttpHandler) {
	router.Handle(http.MethodPut, path, h)
}

func (router *Router) Delete(path string, h HttpHandler) {
	router.Handle(http.MethodDelete, path, h)
}

func (router *Router) Patch(path string, h HttpHandler) {
	router.Handle(http.MethodPatch, path, h)
}

func (router *Router) Options(path string, h HttpHandler) {
	router.Handle(http.MethodOptions, path, h)
}

func (router *Router) Trace(path string, h HttpHandler) {
	router.Handle(http.MethodTrace, path, h)
}

func (router *Router) Connect(path string, h HttpHandler) {
	router.Handle(http.MethodConnect, path, h)
}

func (router *Router) Handle(method string, path string, h HttpHandler) {
	router.configs = append(
		router.configs, config{method, router.prefix + path, router.ms, h},
	)
}

func (router *Router) Group(prefix string, ms *Middlewares, grp GroupCall) {
	newMs := NewMs()
	newMs.Merge(router.ms)
	newMs.Merge(ms)
	router.ms = newMs
	router.prefix += prefix
	grp(router)
	router.ms = NewMs()
	router.prefix = ""
}
