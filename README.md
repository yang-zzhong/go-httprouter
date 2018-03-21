# A Go HTTP Router

```go

router := httprouter.New("/path/to/document/root", []string{"index.html"})

var userList HttpHandler = func(w ResponseWriter, req *Request, _ map[string]string) {
    io.WriteString(w, "user list")
}
var user HttpHandler = func(w ResponseWriter, req *Request, _ map[string]string) {
    io.WriteString(w, "user")
}

router.Group("/api", []Middleware{}, func(router Router) {
    router.Get("/users", UsersList)
    router.Get("/users/:name", User)
})

```
