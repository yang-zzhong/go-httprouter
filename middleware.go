package httprouter

import (
	helpers "github.com/yang-zzhong/go-helpers"
	"net/http"
)

type Middleware interface {
	Before(m http.ResponseWriter, req *Request, p *helpers.P) bool
	After(m http.ResponseWriter, req *Request, p *helpers.P) bool
}

func mergeMiddleware(m1, m2 []Middleware) []Middleware {
	result := []Middleware{}
	for _, mid := range m1 {
		result = append(result, mid)
	}
	for _, mid := range m2 {
		result = append(result, mid)
	}

	return result
}
