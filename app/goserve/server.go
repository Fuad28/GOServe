package server

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/status"
	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

type Server struct {
	Routes         []Route
	Config         Config
	MiddleWares    []HandlerFunc
	AllowedOrigins []string
	Addr           *net.TCPAddr
}

func NewServer(config Config) *Server {
	if config.Port == 0 {
		config.Port = 4221
	}
	if config.TempFileDirectory == "" {
		config.TempFileDirectory = "/tmp/"
	}

	return &Server{
		Config: config,
	}
}

func (s *Server) Address() (*net.TCPAddr, error) {
	localIP, err := getServerIP()
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%d", localIP, s.Config.Port)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Error resolving address: %s", err.Error())
	}

	s.Addr = addr

	return addr, nil
}

func (s *Server) AddRoute(path string, method string, handler HandlerFunc, middlewares []HandlerFunc) (*Route, error) {
	if slices.Contains(REQUEST_METHODS, method) {
		newRoute := NewRoute(path, method, handler, middlewares)
		s.Routes = append(s.Routes, *newRoute)

		return newRoute, nil
	}

	return nil, errors.New("invalid method")
}

func (s *Server) AddMiddleWare(middleware HandlerFunc) {
	s.MiddleWares = append(s.MiddleWares, middleware)
}

func (s *Server) AddAllowedOrigins(addresses []string) {
	s.AllowedOrigins = append(s.AllowedOrigins, addresses...)
}

func (s *Server) GetRoute(req *Request) *Route {
	pathParts := strings.SplitN(req.Path, "?", 2)

	for _, route := range s.Routes {
		isPathMatch, pathParams := matchRoute(pathParts[0], route.path)
		isMethodMatch := req.Method == route.Method

		// For a HEAD request, check that the route's method is HEAD or GET.
		if req.Method == HEAD {
			isMethodMatch = (route.Method == GET) || (route.Method == HEAD)
		}

		// For an OPTIONS request, if the matched route isn't an OPTIONS route, return custom route
		if (isPathMatch) && (req.Method == OPTIONS) && (route.Method != OPTIONS) {
			return DefaultOptionsRoute(s.AllowedOrigins)
		}

		if isPathMatch && isMethodMatch {
			req.PathParams = pathParams

			if len(pathParts) > 1 {
				req.QueryParams = parseQueryParams(pathParts[1])
			} else {
				req.QueryParams = utils.NewKeyValueStore[string, string]()
			}

			return &route
		}
	}

	return nil
}

func (s *Server) HandleRequest(req *Request) IResponse {
	res := NewResponse(req)
	route := s.GetRoute(req)

	if route == nil {
		return res.SetStatus(status.HTTP_404_NOT_FOUND).Send("Path not found.")
	}

	handlerChain := append(append(s.MiddleWares, route.MiddleWares...), route.handler)
	req.HandlerChain = utils.NewQueue[HandlerFunc](handlerChain)

	return req.Next(res)
}

func (s *Server) GET(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, GET, handler, middlewares)
}

func (s *Server) POST(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, POST, handler, middlewares)
}

func (s *Server) PATCH(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, PATCH, handler, middlewares)
}

func (s *Server) PUT(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, PUT, handler, middlewares)
}

func (s *Server) DELETE(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, DELETE, handler, middlewares)
}

func (s *Server) OPTIONS(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, OPTIONS, handler, middlewares)
}

func (s *Server) HEAD(path string, handler HandlerFunc, middlewares ...HandlerFunc) (*Route, error) {
	return s.AddRoute(path, HEAD, handler, append(middlewares, HEADMiddleware))
}

func (s *Server) ServeAndListen() {
	port := s.Config.Port
	l, err := net.Listen("tcp", fmt.Sprint("0.0.0.0:", port))
	defer l.Close()

	if err != nil {
		fmt.Printf("Failed to bind to port %v: %v", port, err.Error())
		os.Exit(1)
	}

	fmt.Println("Server running on port ", port)

	for {
		conn, err := l.Accept()
		defer conn.Close()

		if err != nil {
			fmt.Printf("Error accepting connection: %v", err.Error())
			os.Exit(1)
		}

		clientAddr := conn.RemoteAddr().(*net.TCPAddr)
		serverAddr := s.Addr

		go func() {
			response := NewResponse(nil)
			request := make([]byte, s.Config.MaxRequestSize)
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
			isHead := req.Method == HEAD
			conn.Write(res.GetResponseByte(isHead))
			conn.Close()
		}()
	}
}
