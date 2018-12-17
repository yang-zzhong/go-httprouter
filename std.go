package httprouter

var std *Router

func init() {
	std = NewRouter()
}
func Handler() *Router {
	return std
}
func OnPost(path string, h HttpHandler) {
	std.OnPost(path, h)
}

func OnPut(path string, h HttpHandler) {
	std.OnPut(path, h)
}

func OnDelete(path string, h HttpHandler) {
	std.OnDelete(path, h)
}

func OnGet(path string, h HttpHandler) {
	std.OnGet(path, h)
}

func OnOption(path string, h HttpHandler) {
	std.OnOption(path, h)
}

func OnPatch(path string, h HttpHandler) {
	std.OnPatch(path, h)
}

func OnConnect(path string, h HttpHandler) {
	std.OnConnect(path, h)
}

func Group(prefix string, ms []Middleware, grp GroupCall) {
	std.Group(prefix, ms, grp)
}
