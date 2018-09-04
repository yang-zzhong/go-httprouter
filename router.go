package httprouter

import (
	helper "github.com/yang-zzhong/go-helpers"
	"log"
	"net/http"
)

const (
	Api       = "api"
	PathFile  = "pathfile"
	EntryFile = "entryfile"
)

type HttpHandler func(*ResponseWriter, *Request, *helper.P)
type onFileHandler func(*ResponseWriter, *fileHandler) bool
type GroupCall func(router *Router)
type ResponseHeaderHandler func(*ResponseWriter)
type BeforeExecute func(*ResponseWriter, *Request, *helper.P) bool
type onPanic func(interface{})

type Router struct {
	Tries      []string
	DocRoot    string
	EntryFile  string
	On404      HttpHandler
	BeforeApi  BeforeExecute
	BeforeFile onFileHandler
	OnPanic    onPanic
	configs    []config
	ms         []Middleware
	prefix     string
}

type config struct {
	method string
	path   string
	ms     []Middleware
	call   HttpHandler
}

func onNotFound(w *ResponseWriter, req *Request, _ *helper.P) {
	w.WithStatusCode(http.StatusNotFound)
	w.String("not found")
}

func beforeFile(_ *ResponseWriter, _ *fileHandler) bool {
	return true
}

func beforeApi(_ *ResponseWriter, _ *Request, _ *helper.P) bool {
	return true
}

func NewRouter() *Router {
	router := new(Router)
	router.Tries = []string{Api, PathFile, EntryFile}
	router.OnPanic = func(info interface{}) { log.Print(info) }
	router.DocRoot = "."
	router.EntryFile = "index.html"
	router.BeforeApi = beforeApi
	router.configs = []config{}
	router.ms = []Middleware{}
	router.prefix = ""
	router.On404 = onNotFound
	router.BeforeFile = beforeFile
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.HandleRequest(w, req)
}

func (router *Router) HandleRequest(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			router.OnPanic(err)
		}
	}()
	r := NewResponseWriter(w)
	defer func() {
		log.Printf("%s\t%s\t%v\t%d\t%s", req.Method, req.URL.Path, req.Proto, r.StatusCode, req.RemoteAddr)
		if err := r.Flush(req); err != nil {
			panic(err)
		}
	}()
	if req.Method == http.MethodGet {
		router.try(r, req)
		return
	}
	if router.tryApi(r, req) {
		return
	}
	router.On404(r, &Request{req}, helper.NewP())
}

func (router *Router) try(r *ResponseWriter, req *http.Request) {
	for _, try := range router.Tries {
		switch try {
		case "api":
			if router.tryApi(r, req) {
				return
			}
		case "entryfile":
			if router.tryEntryFile(r, req) {
				return
			}
		case "pathfile":
			if router.tryPathFile(r, req) {
				return
			}
		}
	}
	router.On404(r, &Request{req}, helper.NewP())
}

func (router *Router) tryApi(r *ResponseWriter, req *http.Request) bool {
	methodNotAllowed := false
	for _, conf := range router.configs {
		matched, params := router.Match(conf.method, conf.path, req)
		if !matched {
			continue
		}
		req := &Request{req}
		if !router.BeforeApi(r, req, params) {
			return true
		}
		if req.Method != conf.method {
			methodNotAllowed = true
			continue
		}
		for _, mid := range conf.ms {
			if !mid.Before(r, req, params) {
				return true
			}
			defer mid.After(r, req, params)
		}
		conf.call(r, req, params)

		return true
	}
	if methodNotAllowed {
		r.WithStatusCode(http.StatusMethodNotAllowed)
		return true
	}

	return false
}

func (router *Router) tryEntryFile(r *ResponseWriter, req *http.Request) bool {
	return router.tryFile(r, router.EntryFile)
}

func (router *Router) tryPathFile(r *ResponseWriter, req *http.Request) bool {
	return router.tryFile(r, req.URL.Path)
}

func (router *Router) tryFile(r *ResponseWriter, file string) bool {
	fh := newFileHandler(router.DocRoot)
	available, _ := fh.Available(file)
	if !available {
		return false
	}
	content, cerr := fh.Contents(file)
	if cerr != nil {
		return false
	}
	if router.BeforeFile(r, fh) {
		r.WithHeader("Content-Type", fh.ContentType(file))
		r.WithStatusCode(http.StatusOK)
		r.Write(content)
	}
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

func (router *Router) Group(prefix string, ms []Middleware, grp GroupCall) {
	router.ms = mergeMiddleware(router.ms, ms)
	router.prefix += prefix
	grp(router)
	router.ms = []Middleware{}
	router.prefix = ""
}
