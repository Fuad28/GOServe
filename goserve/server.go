package goserve

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"slices"
	"strings"

	"github.com/Fuad28/GOServe.git/goserve/status"
	"github.com/Fuad28/GOServe.git/goserve/utils"
)

// The server holds it all together and provides methods that handle routes and requests
// "It is it"
type Server struct {

	// Holds all the registered routes
	// The server is the root route.
	// Accessed via routes()
	routes []Route

	// config holds important user-set details for the servers to start.
	// Defaults are set where not provided.
	config Config

	// This holds the middlewares mounted on the server itself, all requests will pass through them.
	// Logging and Localization middlewares can be mounted here.
	middleWares []HandlerFunc

	// These are the allowed orgins that will be permitted if CORSMiddleware is mounted or a pre-flight request is received.
	allowedOrigins []string

	// This is the TCP Address of the server
	addr *net.TCPAddr
}

func (s *Server) Routes() []Route {
	return s.routes
}

func (s *Server) MiddleWares() []HandlerFunc {
	return s.middleWares
}

func (s *Server) AllowedOrigins() []string {
	return s.allowedOrigins
}

// Creates a new server based on config set and returns a pointer to the server instance.
func NewServer(config Config) *Server {
	if config.Port == 0 {
		config.Port = 8000
	}
	if config.MaxRequestSize == 0 {
		config.MaxRequestSize = ONE_MB
	}

	return &Server{
		config: config,
	}
}

// Address is used to obtain the address of the server.
// It's stored in the Addr field and passed to all requests.
func (s *Server) Address() (*net.TCPAddr, error) {
	localIP, err := getServerIP()
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%d", localIP, s.config.Port)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Error resolving address: %s", err.Error())
	}

	s.addr = addr

	return addr, nil
}

// AddRoute is used to register routes on the server.
func (s *Server) AddRoute(path string, method string, handler HandlerFunc, middlewares []HandlerFunc) (*Route, error) {
	if slices.Contains(httpMethods, method) {
		newRoute := NewRoute(path, method, handler, middlewares)
		s.routes = append(s.routes, *newRoute)

		return newRoute, nil
	}

	return nil, errors.New("invalid method")
}

// AddMiddleWare is used to mount middlewares on the server
// e.g server.AddMiddleWare(goserve.CORSMiddleware(route.AllowedOrigins))
func (s *Server) AddMiddleWare(middleware HandlerFunc) {
	s.middleWares = append(s.middleWares, middleware)
}

// AddAllowedOrigins is used to add new trusted origins for CORS after the server has been initialized.
func (s *Server) AddAllowedOrigins(addresses []string) {
	s.allowedOrigins = append(s.allowedOrigins, addresses...)
}

// GetRoute matches requests path and method with registered routes.
// All GET routes also handle HEAD request even when not explictly set.
// All routes handle OPTIONS requests even when not explictly set.
// When handling OPTIONS request when not explictyly set, the route.DefaultOptionsRoute handler is used.

func (s *Server) GetRoute(req *Request) *Route {
	pathParts := strings.SplitN(req.path, "?", 2)

	for _, route := range s.routes {
		isPathMatch, pathParams := matchRoute(pathParts[0], route.path)
		isMethodMatch := req.method == route.method

		// For a HEAD request, check that the route's method is HEAD or GET.
		if req.method == head {
			isMethodMatch = (route.method == get) || (route.method == head)
		}

		// For an OPTIONS request, if the matched route isn't an OPTIONS route, return custom route
		if (isPathMatch) && (req.method == options) && (route.method != options) {
			return DefaultOptionsRoute(s.allowedOrigins)
		}

		if isPathMatch && isMethodMatch {
			req.pathParams = pathParams

			if len(pathParts) > 1 {
				req.queryParams = parseQueryParams(pathParts[1])
			} else {
				req.queryParams = utils.NewKeyValueStore[string, string]()
			}

			return &route
		}
	}

	return nil
}

// HandleRequest processes the requests:
// 1. it initializes the response
// 2. create the HandlerChain and passes the requests into it
// 3. matches registered routes and requests
// 4. handles not found routes
// 5. returns the final response

func (s *Server) HandleRequest(req *Request) IResponse {
	res := NewResponse(req)
	route := s.GetRoute(req)

	if route == nil {
		return res.SetStatus(status.HTTP_404_NOT_FOUND).Send("Path not found.")
	}

	handlerChain := append(append(s.middleWares, route.middleWares...), route.handler)
	req.handlerChain = utils.NewQueue[HandlerFunc](handlerChain)

	return req.Next(res)
}

// GET is shortcut for s.AddRoute(path, get, handler, middlewares)
func (s *Server) GET(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, get, handler, middlewares)
}

// POST is shortcut for s.AddRoute(path, post, handler, middlewares)
func (s *Server) POST(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, post, handler, middlewares)
}

// PATCH is shortcut for s.AddRoute(path, patch, handler, middlewares)
func (s *Server) PATCH(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, patch, handler, middlewares)
}

// PUT is shortcut for s.AddRoute(path, put, handler, middlewares)
func (s *Server) PUT(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, put, handler, middlewares)
}

// DELETE is shortcut for s.AddRoute(path, delete, handler, middlewares)
func (s *Server) DELETE(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, delete, handler, middlewares)
}

// OPTIONS is shortcut for s.AddRoute(path, options, handler, middlewares)
func (s *Server) OPTIONS(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, options, handler, middlewares)
}

// HEAD is shortcut for s.AddRoute(path, head, handler, middlewares)
// The HEADMiddleware is automatically appended to this method to set the body to nil as required by HTTP spec
func (s *Server) HEAD(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, head, handler, append(middlewares, HEADMiddleware))
}

// ServeAndListen is a blocking code that waits for new connections, processes them (asynchronously) and sends responses when done.
// Handles errors that may arise during server start up.
// Handles closing of connections and listner.
func (s *Server) ServeAndListen() {
	port := s.config.Port
	l, err := net.Listen("tcp", fmt.Sprint(":", port))
	defer l.Close()

	if err != nil {
		log.Fatalf("Failed to bind to port %v: %v", port, err.Error())
	}

	log.Println("Server running on port ", port)

	for {
		conn, err := l.Accept()
		defer conn.Close()

		if err != nil {
			log.Fatalf("Error accepting connection: %v", err.Error())
		}

		clientAddr := conn.RemoteAddr().(*net.TCPAddr)
		serverAddr := s.addr

		go func() {
			response := NewResponse(nil)
			request := make([]byte, s.config.MaxRequestSize)
			_, err = conn.Read(request)
			request = bytes.Trim(request, "\x00")

			if err != nil {
				body := fmt.Sprint("Error reading request: ", err.Error())
				response.SetStatus(status.HTTP_400_BAD_REQUEST).Send(body)
				conn.Write(response.GetResponseByte(false))
				return
			}

			req, err := NewRequest(string(request), clientAddr, serverAddr)
			if err != nil {
				body := fmt.Sprint("Error creating request instance: ", err.Error())
				response.SetStatus(status.HTTP_400_BAD_REQUEST).Send(body)
				conn.Write(response.GetResponseByte(false))

				return
			}

			res := s.HandleRequest(req)
			isHead := req.method == head
			conn.Write(res.GetResponseByte(isHead))
			conn.Close()
		}()
	}
}
