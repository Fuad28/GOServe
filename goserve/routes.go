package goserve

type Route struct {
	path        string
	handler     HandlerFunc
	method      string
	middleWares []HandlerFunc
}

func NewRoute(path string, method string, handler HandlerFunc, middlewares []HandlerFunc) *Route {
	return &Route{
		path:        path,
		method:      method,
		handler:     handler,
		middleWares: middlewares,
	}
}

// Handles browsers preflight requests using the CORSMiddleware as an handler
func DefaultOptionsRoute(allowedOrigins []string) *Route {
	return &Route{
		path:        "*",
		method:      options,
		handler:     CORSMiddleware(allowedOrigins),
		middleWares: []HandlerFunc{},
	}
}

func (r *Route) MiddleWares() []HandlerFunc {
	return r.middleWares
}
