package httprouter

import (
	"net/http"
	"strconv"
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
