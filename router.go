package httprouter

import (
	"io"
	"io/ioutil"
	. "net/http"
	"os"
)

type HttpHandler func(ResponseWriter, *Request, map[string]string)
type Middleware func(ResponseWriter, *Request) bool
type group func(router *Router)

type Router struct {
	DocRoot string
	Indexes []string
	configs []config
	mdws    []Middleware
	prefix  string
}

func New(docRoot string, indexes []string) *Router {
	router := new(Router)

	router.DocRoot = docRoot
	router.Indexes = indexes
	router.configs = []config{}
	router.mdws = []Middleware{}
	router.prefix = ""

	return router
}

func (router *Router) ServeHTTP(res ResponseWriter, req *Request) {
	for _, conf := range router.configs {
		matched, params := router.Match(conf.method, conf.path, req)
		if !matched {
			continue
		}
		if req.Method != conf.method {
			res.WriteHeader(StatusMethodNotAllowed)
			return
		}
		conf.Response(res, req, params)
		return
	}
	if req.Method == MethodGet {
		router.defaultRoute(res, req)
		return
	}
	res.WriteHeader(StatusNotFound)
}

func (router *Router) defaultRoute(w ResponseWriter, req *Request) {
	data, err := router.ReadFile(req.URL.Path)
	if err == nil {
		io.WriteString(w, (string)(data))
		return
	}
	pathfile, err := router.IndexFile(req.URL.Path)
	if err != nil {
		w.WriteHeader(StatusNotFound)
		return
	}
	data, err = router.ReadFile(pathfile)
	if err != nil {
		w.WriteHeader(StatusNotFound)
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
	return
	// err = "File Not Found"
}

func (router *Router) Match(method string, p string, req *Request) (matched bool, params map[string]string) {
	matched, params = newPath(p).match(req.URL.Path)
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
		router.configs, config{method, router.prefix + path, router.mdws, h},
	)
}

func (router *Router) Group(prefix string, ms []Middleware, grp group) {
	router.mdws = ms
	router.prefix = prefix
	grp(router)
	router.mdws = []Middleware{}
	router.prefix = ""
}
