package goserve

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"slices"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

type Request struct {
	httpVersion  string
	headers      *utils.KeyValueStore[string, string]
	serverAddr   *net.TCPAddr
	clientAddr   *net.TCPAddr
	host         *url.URL
	origin       *url.URL
	body         any
	method       string
	path         string
	pathParams   *utils.KeyValueStore[string, string]
	queryParams  *utils.KeyValueStore[string, string]
	handlerChain *utils.Queue[HandlerFunc]
	Store        *utils.KeyValueStore[any, any]
}

func NewRequest(req string, clientAddr *net.TCPAddr, serverAddr *net.TCPAddr) (*Request, error) {
	scanner := bufio.NewScanner(strings.NewReader(req))

	if !scanner.Scan() {
		return nil, errors.New("invalid request: missing request line")
	}

	request := Request{
		clientAddr: clientAddr,
		serverAddr: serverAddr,
	}

	// Parse request line
	requestLine := strings.Fields(scanner.Text())
	request.method = strings.ToUpper(requestLine[0])
	request.path = requestLine[1]
	request.httpVersion = requestLine[2]

	// The curl client sends HEAD request using either -I or --head
	if (request.method == "-I") || (request.method == "--head") {
		request.method = "HEAD"
	}

	// Check for invalid request methods
	if !(slices.Contains(httpMethods, request.method)) {
		return nil, errors.New("invalid request: invalid request method")
	}

	// Parse headers
	headers := utils.NewKeyValueStore[string, string]()
	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			break
		}

		headerParts := strings.SplitN(text, ": ", 2)

		if len(headerParts) != 2 {
			return nil, errors.New("invalid request: invalid header")

		} else {
			headers.Set(headerParts[0], headerParts[1])
		}

	}
	request.headers = headers

	// Parse body
	if scanner.Scan() {
		var bodyMap map[string]any
		if err := json.Unmarshal([]byte(scanner.Text()), &bodyMap); err != nil {
			return nil, fmt.Errorf("invalid request: %v", err.Error())
		}
		request.body = bodyMap
	}

	// Parse Host
	hostStr, exists := request.headers.Get("Host")

	if exists {
		host, err := url.Parse("http://" + hostStr)

		if err != nil {
			return nil, errors.New("invalid request: invalid host header")
		}
		request.host = host
	}

	// Parse origin
	originStr, exists := request.headers.Get("Origin")

	if exists {
		origin, err := url.Parse(originStr)

		if err != nil {
			return nil, errors.New("invalid request: invalid origin header")
		}
		request.origin = origin
	}

	return &request, nil
}

func (req *Request) Next(res IResponse) IResponse {
	handler := req.handlerChain.Dequeue().Value
	return handler(req, res)
}

func (req *Request) HTTPVersion() string {
	return req.httpVersion
}

func (req *Request) Headers() *utils.KeyValueStore[string, string] {
	return req.headers
}

func (req *Request) ServerAddr() *net.TCPAddr {
	return req.serverAddr
}

func (req *Request) ClientAddr() *net.TCPAddr {
	return req.clientAddr
}

func (req *Request) Host() *url.URL {
	return req.host
}

func (req *Request) Origin() *url.URL {
	return req.origin
}

func (req *Request) Body() any {
	return req.Body
}

func (req *Request) Method() string {
	return req.method
}

func (req *Request) Path() string {
	return req.path
}

func (req *Request) PathParams() *utils.KeyValueStore[string, string] {
	return req.pathParams
}

func (req *Request) QueryParams() *utils.KeyValueStore[string, string] {
	return req.queryParams
}
