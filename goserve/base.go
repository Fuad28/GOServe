package goserve

import (
	"errors"
	"net"
	"strings"

	"github.com/Fuad28/GOServe.git/goserve/utils"
)

// Signature for route handlers
type HandlerFunc func(*Request, IResponse) IResponse

// Allowed HTTP methods
const (
	get     = "GET"
	post    = "POST"
	put     = "PUT"
	patch   = "PATCH"
	delete  = "DELETE"
	options = "OPTIONS"
	head    = "HEAD"
)

// Array of HTTP allowed httpMethods
var httpMethods = []string{
	get,
	post,
	put,
	patch,
	delete,
	options,
	head,
}

func getServerIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", errors.New("no suitable IP address found")
}

// Utility function used to match request path with registered routes.
// path: the request path gotten from the client.
// route: one of registered routes.
func matchRoute(path string, route string) (bool, *utils.KeyValueStore[string, string]) {
	pathParams := utils.NewKeyValueStore[string, string]()

	pathParts := strings.Split(path, "/")
	routeParts := strings.Split(route, "/")

	if len(pathParts) != len(routeParts) {
		return false, pathParams
	}

	for idx, curPathPart := range pathParts {
		curRoutePart := routeParts[idx]

		if strings.HasPrefix(curRoutePart, ":") {
			paramName := curRoutePart[1:]
			pathParams.Set(paramName, curPathPart)

		} else if curRoutePart != curPathPart {
			return false, pathParams
		}

	}

	return true, pathParams
}

// Utility function to parse query parametes as a key value store.
func parseQueryParams(params string) *utils.KeyValueStore[string, string] {
	queryParams := utils.NewKeyValueStore[string, string]()
	parts := strings.Split(params, "&")

	for _, param := range parts {
		keyValueArr := strings.Split(param, "=")
		queryParams.Set(keyValueArr[0], keyValueArr[1])
	}
	return queryParams
}

// Holds the byte value of 1MB, expected to help with the MaxRequestSize field of the config struct
const ONE_MB = 1024

type JSON map[string]any
