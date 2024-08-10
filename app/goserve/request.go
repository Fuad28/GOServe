package server

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
	HTTPVersion  string
	Headers      *utils.KeyValueStore[string, string]
	ServerAddr   *net.TCPAddr
	ClientAddr   *net.TCPAddr
	Host         *url.URL
	Origin       *url.URL
	Body         any
	Method       string
	Path         string
	PathParams   *utils.KeyValueStore[string, string]
	QueryParams  *utils.KeyValueStore[string, string]
	HandlerChain *utils.Queue[HandlerFunc]
	Store        *utils.KeyValueStore[any, any]
}

func NewRequest(req string, clientAddr *net.TCPAddr, serverAddr *net.TCPAddr) (*Request, error) {
	scanner := bufio.NewScanner(strings.NewReader(req))

	if !scanner.Scan() {
		return nil, errors.New("invalid request: missing request line")
	}

	request := Request{
		ClientAddr: clientAddr,
		ServerAddr: serverAddr,
	}

	// Parse request line
	requestLine := strings.Fields(scanner.Text())
	request.Method = strings.ToUpper(requestLine[0])
	request.Path = requestLine[1]
	request.HTTPVersion = requestLine[2]

	// The curl client sends HEAD request using either -I or --head
	if (request.Method == "-I") || (request.Method == "--head") {
		request.Method = "HEAD"
	}

	// Check for invalid request methods
	if !(slices.Contains(REQUEST_METHODS, request.Method)) {
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
	request.Headers = headers

	// Parse body
	if scanner.Scan() {
		var bodyMap map[string]any
		if err := json.Unmarshal([]byte(scanner.Text()), &bodyMap); err != nil {
			return nil, fmt.Errorf("invalid request: %v", err.Error())
		}
		request.Body = bodyMap
	}

	// Parse Host
	hostStr, exists := request.Headers.Get("Host")

	if exists {
		host, err := url.Parse("http://" + hostStr)

		if err != nil {
			return nil, errors.New("invalid request: invalid host header")
		}
		request.Host = host
	}

	// Parse origin
	originStr, exists := request.Headers.Get("Origin")

	if exists {
		origin, err := url.Parse(originStr)

		if err != nil {
			return nil, errors.New("invalid request: invalid origin header")
		}
		request.Origin = origin
	}

	return &request, nil
}

func (req *Request) Next(res IResponse) IResponse {
	handler := req.HandlerChain.Dequeue().Value
	return handler(req, res)
}
