# A Go HTTP Router

## Feature

1. Support group and middleware
2. With a static server that support front end route
3. Support restful params

```go
import (
    "log"
    "logic" // user app provide
    "net/http"
    hr "github.com/yang-zzhong/go-httprouter"
)

router := hr.NewRouter()

//
// config the try order, here we use the default order
// the router will first match the api, the pathfile based on docroot, the a configed entry file based on docroot
//
router.Tries = []string{hr.API, hr.PATHFILE, hr.ENTRYFILE} 

// only match api
router.Tries = []string{hr.API}

// config docroot
router.DocRoot = "/srv/http/test"

// config api

router.Group("/api", []Mw{}, func(router *Router) {
    router.OnGet("/users", usersList)
    router.OnGet("/users/:user_id", user)
    router.OnPost("/users", createUser)

    router.Group("", []Mw{new(logic.Auth)}, func(router *Router) {
        router.OnPut("/users/:user_id", updateUser)
    });
})

router.OnGet("/hello-world", hello)
log.Fatal(http.ListenAndServe(":8080", router))

var userList HttpHandler = func(w *hr.Response, req *hr.Request) {
    page, _ := req.FormInt("page")
    pageSize, _ := req.FormatInt("page_size")
    w.WithString(logic.UserList(page, pageSize).Json())
}

var user HttpHandler = func(w *hr.Response, req *hr.Request) {
    w.WithString(logic.User(p.Get("user_id")).Json())
}

var createUser HttpHandler = func(w *hr.Response, req *hr.Request) {
    params := map[string]interface{}{
        "name": req.FormValue("name"),
        "account": req.FormValue("account"),
        "extra": req.FormMap("extra"),
    }
    if err := logic.CreateUser(params); err != nil {
        panic(err)
    }
    w.WithString("创建成功")
}

var createUser HttpHandler = func(w *hr.Response, req *hr.Request) {
    params := map[string]interface{}{
        "name": req.FormValue("name"),
        "account": req.FormValue("account"),
        "extra": req.FormMap("extra"),
    }
    if err := logic.UpdateUser(req.Bag.Get("user_id"), params); err != nil {
        panic(err)
    }

    w.WithString("更新成功")
}

var hello HttpHandler = func(w *hr.Response, _ *hr.Request) {
    w.WithString("hello world!!!")
}

```
