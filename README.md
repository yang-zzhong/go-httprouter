# A Go HTTP Router

## Feature

1. Support group and middleware
2. With a static server that support front end route
3. Support restful params

```go
import (
    . "net/http"
    "io"
    helper "github.com/yang-zzhong/go-helpers"
    httprouter "github.com/yang-zzhong/go-httprouter"
)

router := httprouter.NewRouter()

//
// config the try order, here we use the default order
// the router will first match the api, the pathfile based on docroot, the a configed entry file based on docroot
//
router.Tries = []string{httprouter.Api, httprouter.PathFile, httprouter.EntryFile} 

// only match api
router.Tries = []string{httprooter.Api}

// config docroot
router.DocRoot = "/srv/http/test"

// config api
var userList HttpHandler = func(w ResponseWriter, req *Request, _ *helper.P) {
    io.WriteString(w, "user list")
}
var user HttpHandler = func(w ResponseWriter, req *Request, _ *helper.P) {
    io.WriteString(w, "user")
}

var hello HttpHandler = func(w ResponseWriter, req *Request, p *helper.P) {
    io.WriteString(w, "hello " + p.Get("world"))
}

router.Group("/api", NewMs(), func(router Router) {
    router.Get("/users", UsersList)
    router.Get("/users/:name", User)
})
router.Get("/hello/:world", hello)

```
