package goserve

import (
	"slices"

	"github.com/Fuad28/GOServe.git/goserve/status"
)

// Middlewares provide the functionality to intercept and modify the request and response at any point.
// Middlewares are functions that have the HandlerFunc signature.
// Middleware functions must always call the next middleware (req.Next(res)) to pass control to the next middleware.
// You can access the request/ response at any point in the request lifecycle as recursion is used to achieve this e.g check the HEADMiddleware.
// It can be mounted directly on the main router (i.e the server) i.e route.AddMiddleWare(goserve.CORSMiddleware(route.AllowedOrigins))
// It can be mounted on indiviual route at the point of registering: route.GET("/tasks", tasks, middleware1, middleware2)
// In the grand scheme of things, the request-response cycle is simply an entire middlewares chain.
// The CORSMiddleware is provided to allow CORS implementation. It's not mounted by default
// The HEADMiddleware is used to majorly to set the response body to null when handling OPTIONS requests.
// You can use it if you find other applications fot it.

// If the CORSMiddleware is mounted, it intercepts all requests and check if the clientAddr is in the server's AllowedOrigins array.
// If it's not present, we immediately return a 403 Forbidden response.
// If it's present and it's an OPTIONS request (i.e preflight), we return a 200 OK response with neccessary CORS headers set.
// If it's present and it's not a preflight, we pass control to the next middleware.
func CORSMiddleware(allowedOrigins []string) HandlerFunc {
	return func(req *Request, res IResponse) IResponse {
		isAllowed := false
		origin := ""

		if (req.origin != nil) && (req.host != nil) {
			origin = req.origin.String()
			isSameOrigin := req.origin.Hostname() == req.host.Hostname()
			isAllowedOrigin := slices.Contains(allowedOrigins, req.origin.String())
			isAllowed = isSameOrigin || isAllowedOrigin
		}

		if !isAllowed {
			return res.SetStatus(status.HTTP_403_FORBIDDEN).Send("Forbidden.")
		}

		res.SetHeader("Access-Control-Allow-Origin", origin)
		res.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		res.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if req.method == options {
			return res.SetStatus(status.HTTP_200_OK).Send(nil)
		}

		return req.Next(res)
	}
}

// The HEADMiddleware is appended last in the list of the middlewares for routes with Head method.
// It sets the body to nil as required for HEAD request method.
func HEADMiddleware(req *Request, res IResponse) IResponse {
	res = req.Next(res)
	res.Send(nil)

	return res
}
