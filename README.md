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
var userList HttpHandler = func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *helper.P) {
    page := req.FormInt("page")
    pageSize := req.FormatInt("page_size")
    w.WriteJson(logic.UserList(page, pageSize))
}

var user HttpHandler = func(w *httprouter.ResponseWriter, _ *httprouter.Request, p *helper.P) {
    w.WriteJson(logic.User(p.Get("user_id")))
}

var createUser HttpHandler = func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *helpers.P) {
    params := map[string]interface{}{
        "name": req.FormValue("name"),
        "account": req.FormValue("account"),
        "extra": req.FormMap("extra"),
    }
    if err := logic.CreateUser(params); err != nil {
        panic(err)
    }
    w.WriteString("创建成功")
}

var createUser HttpHandler = func(w *httprouter.ResponseWriter, req *httprouter.Request, p *helpers.P) {
    params := map[string]interface{}{
        "name": req.FormValue("name"),
        "account": req.FormValue("account"),
        "extra": req.FormMap("extra"),
    }
    if err := logic.UpdateUser(p.Get("user_id"), params); err != nil {
        panic(err)
    }

    w.WriteString("更新成功")
}

var hello HttpHandler = func(w *httprouter.ResponseWriter, _ *httprouter.Request, _ *helper.P) {
    w.WriteString("hello world!!!")
}

router.Group("/api", []Middleware{}, func(router *Router) {
    router.OnGet("/users", usersList)
    router.OnGet("/users/:user_id", user)
    router.OnPost("/users", createUser)

    router.Group("", []Middleware{new(logic.Auth)}, func(router *Router) {
        router.OnPut("/users/:user_id", updateUser)
    });
})

router.OnGet("/hello-world", hello)

log.Fatal(http.ListenAndServe(":8080", router))

```
