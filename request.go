package httprouter

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Bagt struct {
	body map[string]interface{}
}

func NewBagt() *Bagt {
	return &Bagt{make(map[string]interface{})}
}

func (p *Bagt) Set(k string, v interface{}) {
	p.body[k] = v
}

func (p *Bagt) Get(k string) interface{} {
	return p.body[k]
}

func (p *Bagt) Del(k string) {
	if p.Exist(k) {
		delete(p.body, k)
	}
}

func (p *Bagt) Each(handle func(k string, v interface{}) bool) bool {
	for k, v := range p.body {
		if !handle(k, v) {
			return false
		}
	}
	return true
}

func (p *Bagt) Exist(k string) bool {
	_, ok := p.body[k]
	return ok
}

// wrap http.Request, and provide some useful functions
type Request struct {
	Bag *Bagt
	*http.Request
}

// read form field as int64, if you need other int type, use type convert
func (req *Request) FormInt(fieldname string) (r int64, e error) {
	val := req.FormValue(fieldname)
	if val == "" {
		return 0, errors.New("field not found")
	}
	r, e = strconv.ParseInt(val, 10, 64)
	return
}

// read form field as uint64, if you need other uint type, use type convert
func (req *Request) FormUint(fieldname string) (r uint64, e error) {
	val := req.FormValue(fieldname)
	if val == "" {
		return 0, errors.New("field not found")
	}
	r, e = strconv.ParseUint(val, 10, 64)
	return
}

// read form field as float, if you need other float type, use type convert
func (req *Request) FormFloat(fieldname string) (r float64, e error) {
	val := req.FormValue(fieldname)
	if val == "" {
		return 0.0, errors.New("field not found")
	}
	r, e = strconv.ParseFloat(val, 64)
	return
}

// read form field as bool, "false", 0, "" will be recognised as false, others true
func (req *Request) FormBool(fieldName string) bool {
	val := req.FormValue(fieldName)

	return !(val == "" || val == "0" || val == "false")
}

// read form field as string slice, key[\d+] will be recognised item
func (req *Request) FormSlice(fieldname string) []string {
	r := strings.Split(req.FormValue(fieldname), ",")
	req.ParseForm()
	for key, val := range req.Form {
		var matched bool
		matched, _ = regexp.MatchString("^"+fieldname+"\\[\\d*\\]$", key)
		if !matched {
			continue
		}
		r = append(r, val...)
	}
	res := []string{}
	for _, val := range r {
		if val != "" {
			res = append(res, val)
		}
	}

	return res
}

// read form field as string slice, key[\w+] will be recognised item
func (req *Request) FormMap(fieldname string) map[string]string {
	result := make(map[string]string)
	var field string
	var reg *regexp.Regexp
	var keyReg *regexp.Regexp
	reg, _ = regexp.Compile("^" + fieldname + "\\[\\w+\\]$")
	keyReg, _ = regexp.Compile("\\[\\w+\\]")
	req.ParseForm()
	for key, val := range req.Form {
		field = keyReg.FindString(reg.FindString(key))
		if field != "" {
			result[field[1:len(field)-1]] = val[0]
		}
	}

	return result
}
