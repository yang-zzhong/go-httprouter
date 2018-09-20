package httprouter

import (
	helpers "github.com/yang-zzhong/go-helpers"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// wrap http.Request, and provide some useful functions
type Request struct {
	Bag *helpers.P
	*http.Request
}

// read form field as int64, if you need other int type, use type convert
func (req *Request) FormInt(fieldName string) int64 {
	val := req.FormValue(fieldName)
	if val == "" {
		return 0
	}
	var result int64
	var err error
	if result, err = strconv.ParseInt(val, 10, 64); err != nil {
		return 0
	}
	return result
}

// read form field as uint64, if you need other uint type, use type convert
func (req *Request) FormUint(fieldName string) uint64 {
	val := req.FormValue(fieldName)
	if val == "" {
		return 0
	}
	var result uint64
	var err error
	if result, err = strconv.ParseUint(val, 10, 64); err != nil {
		return 0
	}

	return result
}

// read form field as float, if you need other float type, use type convert
func (req *Request) FormFloat(fieldName string) float64 {
	val := req.FormValue(fieldName)
	if val == "" {
		return 0.0
	}
	var result float64
	var err error
	if result, err = strconv.ParseFloat(val, 64); err != nil {
		return 0.0
	}
	return result
}

// read form field as bool, "false", 0, "" will be recognised as false, others true
func (req *Request) FormBool(fieldName string) bool {
	val := req.FormValue(fieldName)

	return !(val == "" || val == "0" || val == "false")
}

// read form field as string slice, key[\d+] will be recognised item
func (req *Request) FormSlice(fieldName string) []string {
	result := strings.Split(req.FormValue(fieldName), ",")
	req.ParseForm()
	for key, val := range req.Form {
		var matched bool
		var err error
		if matched, err = regexp.MatchString("^"+fieldName+"\\[\\d*\\]$", key); err != nil {
			log.Println(err)
			continue
		}
		if !matched {
			continue
		}
		result = append(result, val...)
	}
	res := []string{}
	for _, val := range result {
		if val != "" {
			res = append(res, val)
		}
	}

	return res
}

// read form field as string slice, key[\w+] will be recognised item
func (req *Request) FormMap(fieldName string) map[string]string {
	result := make(map[string]string)
	var err error
	var field string
	var reg *regexp.Regexp
	var keyReg *regexp.Regexp
	if reg, err = regexp.Compile("^" + fieldName + "\\[\\w+\\]$"); err != nil {
		log.Println(err)
		return result
	}
	if keyReg, err = regexp.Compile("\\[\\w+\\]"); err != nil {
		log.Println(err)
		return result
	}
	req.ParseForm()
	for key, val := range req.Form {
		field = keyReg.FindString(reg.FindString(key))
		if field != "" {
			result[field[1:len(field)-1]] = val[0]
		}
	}

	return result
}
