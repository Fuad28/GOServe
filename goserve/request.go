package goserve

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"github.com/Fuad28/GOServe.git/goserve/utils"
)

type Request struct {
	// Represents the HTTPVersion of the request, this is used in constructing the response.
	// Accessed via HTTPVersion()
	httpVersion string

	// Holds the values of the request headers
	// Accessed via Headers()
	headers *utils.KeyValueStore[string, string]

	// Is the TCP address of the server.
	//This isn't used internally but seen as a valuable data to have.
	// Accessed via ServerAddr()
	serverAddr *net.TCPAddr

	// The TCP address of the client making the request.
	// This isn't used internally but seen as a valuable data to have.
	// Accessed via ClientAddr()
	clientAddr *net.TCPAddr

	// host is the domain the server is running on and is expected to be set by the client.
	// Used in evaluating same origin requests if CORS middleware is mounted or in handling a pre-flight (OPTIONS) request.
	// Accessed via Host()
	host *url.URL

	// origin is the domain of the client is expected to be set by the client.
	// Used in evaluating if the server is able to share with the client when the CORS middleware mounted or in handling a pre-flight (OPTIONS) request.
	// Accessed via Origin()
	origin *url.URL

	// Holds the body of the request which is expected to be valid JSON serializatble.
	// Accessed via Body()
	body []byte

	// Request method
	// Accessed via Method()
	method string

	// Request target/endpoint/path
	// Accessed via Path()
	path string

	// uses the *utils.KeyValueStore[string, string] data structure to hold path parameters found in the request path.
	// Accessed via PathParams()
	pathParams *utils.KeyValueStore[string, string]

	// uses the *utils.KeyValueStore[string, string] data structure to hold query parameters found in the request path.
	// Accessed via QueryParams()
	queryParams *utils.KeyValueStore[string, string]

	// uses the *utils.Queue[HandlerFunc] to hold the entire handlers chain for the request.
	// While the request is being handled, the middlewares and handler are put in a queue to preserve order and allow for efficient retrieval.
	handlerChain *utils.Queue[HandlerFunc]

	// An empty Store of type *utils.KeyValueStore[string, string] is kept on all requests.
	// Allows for sotring and passing data throughout the request-response cycle.
	Store *utils.KeyValueStore[any, any]
}

func NewRequest(req string, clientAddr *net.TCPAddr, serverAddr *net.TCPAddr) (*Request, error) {
	scanner := bufio.NewScanner(strings.NewReader(req))

	if !scanner.Scan() {
		return nil, errors.New("invalid request: missing request line")
	}

	request := Request{
		clientAddr: clientAddr,
		serverAddr: serverAddr,
		Store:      utils.NewKeyValueStore[any, any](),
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

		// Validate body is a valid JSON
		if body, err := json.Marshal(scanner.Text()); err != nil {
			return nil, fmt.Errorf("invalid request: %v", err.Error())

		} else {
			request.body = body
		}
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

func (req *Request) Body(v any) error {

	if reflect.TypeOf(v).Kind() != reflect.Pointer {
		log.Fatal("v must be a point")
	}

	return json.Unmarshal(req.body, v)
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
