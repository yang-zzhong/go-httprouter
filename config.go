package httprouter

import (
	. "net/http"
)

type config struct {
	method string
	path   string
	mdws   []Middleware
	call   HttpHandler
}

func (conf *config) Response(w ResponseWriter, req *Request, params map[string]string) {
	for _, middleware := range conf.mdws {
		if middleware(w, req) {
			continue
		}
		return
	}

	conf.call(w, req, params)
}
