package server

type Route struct {
	path        string
	handler     HandlerFunc
	Method      string
	MiddleWares []HandlerFunc
}

func NewRoute(path string, method string, handler HandlerFunc, middlewares []HandlerFunc) *Route {
	return &Route{
		path:        path,
		Method:      method,
		handler:     handler,
		MiddleWares: middlewares,
	}
}

// Handles browsers preflight requests using the CORSMiddleware as an handler
func DefaultOptionsRoute(allowedOrigins []string) *Route {
	return &Route{
		path:        "*",
		Method:      OPTIONS,
		handler:     CORSMiddleware(allowedOrigins),
		MiddleWares: []HandlerFunc{},
	}
}
