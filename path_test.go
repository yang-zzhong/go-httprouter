package httprouter

import (
	"testing"
)

func TestSimple(t *testing.T) {

	p := &path{"/users"}
	matched, _ := p.match("/users")
	if !matched {
		t.Error("/users Error")
	}
}

func TestParams(t *testing.T) {
	p := &path{"/users/:name"}
	matched, params := p.match("/users/yang-zzhong")
	if !matched {
		t.Error("/users/:name Not Matched")
	}
	if params.Get("name") != "yang-zzhong" {
		t.Error("/users/:name Params Not Got")
	}
}

func TestComplex(t *testing.T) {
	p := &path{"/users/:name/articles"}
	matched, params := p.match("/users/yang-zzhong/articles")
	if !matched {
		t.Error("/users/:name/articles Not Matched")
	}
	if params.Get("name") != "yang-zzhong" {
		t.Error("/users/:name/articles Params Not Got")
	}
}

func TestAllParams(t *testing.T) {
	p := &path{"/:name/:article-name"}
	matched, params := p.match("/yang-zzhong/Me+And+My+Broken+Heart")
	if !matched {
		t.Error("/:name/:article-name Not Matched")
	}
	para := make(map[string]string)
	para["name"] = "yang-zzhong"
	para["article-name"] = "Me+And+My+Broken+Heart"
	params.Each(func(key string, value interface{}) bool {
		if para[key] != value {
			t.Error("/:name/:article-name Params Not Got")
			return false
		}

		return true
	})
}
