# A Go HTTP Router

## Feature

1. Support group and middleware
2. A default handler to find file from document root
3. Support restful params

```go
import (
    . "net/http"
    "io"
    helper "github.com/yang-zzhong/go-helpers"
    httprouter "github.com/yang-zzhong/go-httprouter"
)

router := httprouter.CreateRouter("/path/to/document/root", []string{"index.html"})

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
