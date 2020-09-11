// the package provide a simple clean http route on server side
package httprouter

import (
	restfulpath "github.com/yang-zzhong/go-restfulpath"
	"net/http"
	"os"
	. "path"
)

const (
	API = iota
	PATHFILE
	ENTRYFILE
)

// http handler type
type HttpHandler func(*Response, *Request)

// router as file server, when output file, execute the callback. here is the type
type onFileHandler func(*Response, *http.Request, string) bool

// group call type
type GroupCall func(router *Router)

// when uri match, the callback will be executed. warning that, when different method, uri possibly match many times
type BeforeExecute func(http.ResponseWriter, *http.Request) bool

type Router struct {
	Tries           []int
	DocRoot         string
	EntryFile       string
	BeforePathFile  onFileHandler
	BeforeEntryFile onFileHandler
	configs         []config
	ms              []Mw
	prefix          string
}

type config struct {
	method string
	path   string
	ms     []Mw
	call   HttpHandler
}

func beforeFile(_ *Response, _ *http.Request, _ string) bool {
	return true
}

// new router
func NewRouter() *Router {
	router := new(Router)
	router.Tries = []int{API, PATHFILE, ENTRYFILE}
	router.DocRoot = "."
	router.EntryFile = "index.html"
	router.configs = []config{}
	router.ms = []Mw{}
	router.prefix = ""
	router.BeforePathFile = beforeFile
	router.BeforeEntryFile = beforeFile
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := router.HandleRequest(w, req)
	r.Flush(req)
}

func (router *Router) HandleRequest(w http.ResponseWriter, req *http.Request) *Response {
	r := NewResponse(w)
	if req.Method == http.MethodGet {
		router.try(r, req)
	} else {
		router.tryApi(r, req)
	}
	return r
}

func (router *Router) try(r *Response, req *http.Request) {
	for _, try := range router.Tries {
		switch try {
		case API:
			if router.tryApi(r, req) {
				return
			}
		case PATHFILE:
			if router.tryPathFile(r, req) {
				return
			}
		case ENTRYFILE:
			if router.tryEntryFile(r, req) {
				return
			}
		}
	}
}

func (router *Router) tryApi(r *Response, req *http.Request) bool {
	methodNotAllowed := false
	for _, conf := range router.configs {
		matched, params := restfulpath.NewPath(conf.path).Match(req.URL.Path)
		if !matched {
			continue
		}
		if req.Method != conf.method {
			methodNotAllowed = true
			continue
		}
		bag := NewBagt()
		for k, v := range params {
			bag.Set(k, v)
		}
		wreq := &Request{bag, req}
		for _, mid := range conf.ms {
			if !mid.Before(r, wreq) {
				return true
			}
			defer mid.After(r, wreq)
		}
		conf.call(r, wreq)

		return true
	}
	if methodNotAllowed {
		r.WithStatus(http.StatusMethodNotAllowed)
		return true
	}

	return false
}

func (router *Router) tryEntryFile(r *Response, req *http.Request) bool {
	return router.tryFile(r, req, router.EntryFile, router.BeforeEntryFile)
}

func (router *Router) tryPathFile(r *Response, req *http.Request) bool {
	return router.tryFile(r, req, req.URL.Path, router.BeforePathFile)
}

func (router *Router) tryFile(r *Response, req *http.Request, file string, beforeFile onFileHandler) bool {
	pathfile := Join(router.DocRoot, file)
	if stat, err := os.Stat(pathfile); err != nil {
		if os.IsNotExist(err) {
			r.WithStatus(404).WithString("File Not Found")
			return false
		}
	} else if stat.IsDir() {
		r.WithStatus(404).WithString("File Not Found")
		return false
	}
	r.WithStatus(200)
	if beforeFile(r, req, pathfile) {
		r.WithFile(pathfile)
	}
	return true
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
func (router *Router) Group(prefix string, ms []Mw, grp GroupCall) {
	router.ms = mergeMiddleware(router.ms, ms)
	router.prefix += prefix
	grp(router)
	router.ms = []Mw{}
	router.prefix = ""
}
