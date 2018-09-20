package httprouter

type Middleware interface {
	Before(m *ResponseWriter, req *Request) bool
	After(m *ResponseWriter, req *Request) bool
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
