package httprouter

import (
	helpers "github.com/yang-zzhong/go-helpers"
	"net/http"
)

type Middleware func(http.ResponseWriter, *Request, *helpers.P) bool

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

func (ms *Middlewares) Exec(w http.ResponseWriter, req *Request, p *helpers.P) bool {
	for _, middleware := range ms.mdws {
		if !middleware(w, req, p) {
			return false
		}
	}

	return true
}
