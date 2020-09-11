package httprouter

type Mw interface {
	Before(m *Response, req *Request) bool
	After(m *Response, req *Request) bool
}

func mergeMiddleware(m1, m2 []Mw) []Mw {
	result := []Mw{}
	for _, mid := range m1 {
		result = append(result, mid)
	}
	for _, mid := range m2 {
		result = append(result, mid)
	}

	return result
}
