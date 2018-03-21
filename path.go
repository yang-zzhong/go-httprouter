package httprouter

import (
	. "bytes"
	helper "github.com/yang-zzhong/go-helpers"
	"regexp"
)

// /users
// /users/:name
type path struct {
	p string
}

func newPath(p string) *path {
	pa := new(path)
	pa.p = p

	return pa
}

func (p *path) match(t string) (matched bool, params *helper.P) {
	params = helper.NewP()
	pa := Split(([]byte)(p.p), []byte{'/'})
	lenpa := len(pa)
	ta := Split(([]byte)(t), []byte{'/'})
	lenta := len(ta)
	if lenpa != lenta {
		matched = false
		return
	}
	for i, ip := range pa {
		ips := (string)(ip)
		its := (string)(ta[i])
		m, _ := regexp.Match("^:", ip)
		if m {
			params.Set(ips[1:len(ips)], its)
			continue
		}
		if ips != its {
			matched = false
			return
		}
	}
	matched = true

	return
}
