// the package provide a simple clean http route on server side
package httprouter

import (
	helper "github.com/yang-zzhong/go-helpers"
	"runtime/debug"
	"log"
	"net/http"
	"os"
	. "path"
)

const (
	Api       = "api"
	PathFile  = "pathfile"
	EntryFile = "entryfile"
)

// http handler type
type HttpHandler func(*ResponseWriter, *Request)

// router as file server, when output file, execute the callback. here is the type
type onFileHandler func(*ResponseWriter, *http.Request, string) bool

// group call type
type GroupCall func(router *Router)

// when uri match, the callback will be executed. warning that, when different method, uri possibly match many times
type BeforeExecute func(*ResponseWriter, *Request) bool

// any panic will cause the callback execute
type onPanic func(interface{}, *ResponseWriter, *http.Request)

type Router struct {
	Tries           []string
	DocRoot         string
	EntryFile       string
	On404           HttpHandler
	BeforeApi       BeforeExecute
	BeforePathFile  onFileHandler
	BeforeEntryFile onFileHandler
	OnPanic         onPanic
	configs         []config
	ms              []Middleware
	prefix          string
}

type config struct {
	method string
	path   string
	ms     []Middleware
	call   HttpHandler
}

func onNotFound(w *ResponseWriter, req *Request) {
	w.WithStatusCode(http.StatusNotFound)
	w.String("not found")
}

func beforeFile(_ *ResponseWriter, _ *http.Request, _ string) bool {
	return true
}

func beforeApi(_ *ResponseWriter, _ *Request) bool {
	return true
}

// new router
func NewRouter() *Router {
	router := new(Router)
	router.Tries = []string{Api, PathFile, EntryFile}
	router.OnPanic = func(info interface{}, w *ResponseWriter, req *http.Request) {
		w.WithStatusCode(500)
		w.String("Server Unknown Error")
		debug.PrintStack()
	}
	router.DocRoot = "."
	router.EntryFile = "index.html"
	router.BeforeApi = beforeApi
	router.configs = []config{}
	router.ms = []Middleware{}
	router.prefix = ""
	router.On404 = onNotFound
	router.BeforePathFile = beforeFile
	router.BeforeEntryFile = beforeFile
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.HandleRequest(w, req)
}

func (router *Router) HandleRequest(w http.ResponseWriter, req *http.Request) {
	r := NewResponseWriter(w)
	defer func() {
		log.Printf("%s\t%s\t%v\t%d\t%s", req.Method, req.URL.Path, req.Proto, r.StatusCode, req.RemoteAddr)
		if err := r.Flush(req); err != nil {
			panic(err)
		}
	}()
	defer func() {
		if e := recover(); e != nil {
			router.OnPanic(e, r, req)
		}
	}()
	if req.Method == http.MethodGet {
		router.try(r, req)
		return
	}
	if router.tryApi(r, req) {
		return
	}
	router.On404(r, &Request{helper.NewP(), req})
}

func (router *Router) try(r *ResponseWriter, req *http.Request) {
	for _, try := range router.Tries {
		switch try {
		case "api":
			if router.tryApi(r, req) {
				return
			}
		case "pathfile":
			if router.tryPathFile(r, req) {
				return
			}
		case "entryfile":
			if router.tryEntryFile(r, req) {
				return
			}
		}
	}
	router.On404(r, &Request{helper.NewP(), req})
}

func (router *Router) tryApi(r *ResponseWriter, req *http.Request) bool {
	methodNotAllowed := false
	for _, conf := range router.configs {
		matched, params := router.Match(conf.method, conf.path, req)
		if !matched {
			continue
		}
		req := &Request{params, req}
		if !router.BeforeApi(r, req) {
			return true
		}
		if req.Method != conf.method {
			methodNotAllowed = true
			continue
		}
		for _, mid := range conf.ms {
			if !mid.Before(r, req) {
				return true
			}
			defer mid.After(r, req)
		}
		conf.call(r, req)

		return true
	}
	if methodNotAllowed {
		r.WithStatusCode(http.StatusMethodNotAllowed)
		return true
	}

	return false
}

func (router *Router) tryEntryFile(r *ResponseWriter, req *http.Request) bool {
	return router.tryFile(r, req, router.EntryFile, router.BeforeEntryFile)
}

func (router *Router) tryPathFile(r *ResponseWriter, req *http.Request) bool {
	return router.tryFile(r, req, req.URL.Path, router.BeforePathFile)
}

func (router *Router) tryFile(r *ResponseWriter, req *http.Request, file string, beforeFile onFileHandler) bool {
	pathfile := Join(router.DocRoot, file)
	if stat, err := os.Stat(pathfile); err != nil {
		if os.IsNotExist(err) {
			r.WithStatusCode(404).String("File Not Found")
			return false
		}
	} else if stat.IsDir() {
		r.WithStatusCode(404).String("File Not Found")
		return false
	}
	r.WithStatusCode(200)
	if beforeFile(r, req, pathfile) {
		r.WriteFile(pathfile)
	}
	return true
}

func (router *Router) Match(method string, path string, req *http.Request) (m bool, p *helper.P) {
	m, p = newPath(path).match(req.URL.Path)
	return
}

// on get uri
func (router *Router) OnGet(path string, h HttpHandler) {
	router.Get(path, h)
}

// on post uri
func (router *Router) OnPost(path string, h HttpHandler) {
	router.Post(path, h)
}

// on put uri
func (router *Router) OnPut(path string, h HttpHandler) {
	router.Put(path, h)
}

// on delete uri
func (router *Router) OnDelete(path string, h HttpHandler) {
	router.Delete(path, h)
}

// on patch uri
func (router *Router) OnPatch(path string, h HttpHandler) {
	router.Patch(path, h)
}

// on connect uri
func (router *Router) OnConnect(path string, h HttpHandler) {
	router.Connect(path, h)
}

// on option uri
func (router *Router) OnOption(path string, h HttpHandler) {
	router.Option(path, h)
}

// on trace uri
func (router *Router) OnTrace(path string, h HttpHandler) {
	router.OnTrace(path, h)
}

// legacy on get uri, abandon future
func (router *Router) Get(path string, h HttpHandler) {
	router.Handle(http.MethodGet, path, h)
}

// legacy on post uri, abandon future
func (router *Router) Post(path string, h HttpHandler) {
	router.Handle(http.MethodPost, path, h)
}

// legacy on put uri, abandon future
func (router *Router) Put(path string, h HttpHandler) {
	router.Handle(http.MethodPut, path, h)
}

// legacy on delete uri, abandon future
func (router *Router) Delete(path string, h HttpHandler) {
	router.Handle(http.MethodDelete, path, h)
}

// legacy on patch uri, abandon future
func (router *Router) Patch(path string, h HttpHandler) {
	router.Handle(http.MethodPatch, path, h)
}

// legacy on options uri, abandon future
func (router *Router) Option(path string, h HttpHandler) {
	router.Handle(http.MethodOptions, path, h)
}

// legacy on trace uri, abandon future
func (router *Router) Trace(path string, h HttpHandler) {
	router.Handle(http.MethodTrace, path, h)
}

// legacy on connect uri, abandon future
func (router *Router) Connect(path string, h HttpHandler) {
	router.Handle(http.MethodConnect, path, h)
}

// handle a request
func (router *Router) Handle(method string, path string, h HttpHandler) {
	router.configs = append(
		router.configs, config{method, router.prefix + path, router.ms, h},
	)
}

// add prefix, middleware for a bunch of request
func (router *Router) Group(prefix string, ms []Middleware, grp GroupCall) {
	router.ms = mergeMiddleware(router.ms, ms)
	router.prefix += prefix
	grp(router)
	router.ms = []Middleware{}
	router.prefix = ""
}
