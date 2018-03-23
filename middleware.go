package httprouter

import (
	. "net/http"
)

type Middleware func(ResponseWriter, *Request) bool

type Middlewares struct {
	mdws []Middleware
}

func NewMs() *Middlewares {
	ms := new(Middlewares)
	ms.mdws = []Middleware{}
	return ms
}

func (ms *Middlewares) Append(md Middleware) *Middlewares {
	ms.mdws = append(ms.mdws, md)

	return ms
}

func (ms *Middlewares) Merge(mms *Middlewares) *Middlewares {
	for _, m := range mms.All() {
		ms.Append(m)
	}
	return ms
}

func (ms *Middlewares) All() []Middleware {
	return ms.mdws
}

func (ms *Middlewares) Exec(w ResponseWriter, req *Request) bool {
	for _, middleware := range ms.mdws {
		if !middleware(w, req) {
			return false
		}
	}

	return true
}
