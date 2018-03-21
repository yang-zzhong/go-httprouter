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

func (ms *Middlewares) Append(md Middleware) {
	ms.mdws = append(ms.mdws, md)
}

func (ms *Middlewares) Exec(w ResponseWriter, req *Request) bool {
	for _, middleware := range ms.mdws {
		if middleware(w, req) {
			continue
		}
		return false
	}

	return true
}
