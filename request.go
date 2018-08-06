package httprouter

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Request struct {
	*http.Request
}

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

func (req *Request) FormBool(fieldName string) bool {
	val := req.FormValue(fieldName)

	return !(val == "" || val == "0" || val == "false")
}

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
