package goserve

import (
	"slices"

	"github.com/Fuad28/GOServe.git/status"
)

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
