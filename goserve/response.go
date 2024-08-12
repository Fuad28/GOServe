package goserve

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Fuad28/GOServe.git/goserve/status"
	"github.com/Fuad28/GOServe.git/goserve/utils"
)

// The IResponse defines the interface for a Response
type IResponse interface {

	// HTTPVersion returns a value representing the HTTP version used in the sending the response.
	HTTPVersion() string

	// StatusCode returns the status code of the response. It defaults to 200 if not set.
	StatusCode() int

	// This is used to set the response status code e.g res.SetStatus(status.HTTP_200_OK)
	SetStatus(int) IResponse

	// Sets an header. Will override an header if key exists.
	SetHeader(key string, value string) IResponse

	// Gives access to the response headers as *utils.KeyValueStore[string, string]
	// e.g res.Headers().Get("Content-Type")
	Headers() *utils.KeyValueStore[string, string]

	// Allows you to access the response body
	Body() any

	// Sets the response body.
	// Note that: This method only sends the response body and doesn't send the response, you have to return the 'res' in your controller.
	Send(any) IResponse

	// This gives the byte array representation of the response body.
	// This is invoked in the request-response cycle after a response is ready.
	// It accepts an isHead bool to know whether to set request body or not.
	GetResponseByte(bool) []byte
}

// Response is a type that holds response data.
// It's an implementation of the IResponse interface.
type Response struct {
	// Tihs is the protocol over which the response is sent.
	// It is extracted from the request.
	// Accessed via HTTPVersion()
	httpVersion string

	// The is the status code for the response. Defaults to status.HTTP_200_OK
	// Accessed via Status()
	statusCode int

	// Holds the values of the response headers set.
	// The Content-Type and Content-Length headers are set by default just before response is sent
	// Accessed via Headers()
	headers *utils.KeyValueStore[string, string]

	// Holds the body of the reposne which is expected to be valid JSON serializatble.
	// Accessed via Body()
	body any
}

func NewResponse(req *Request) *Response {
	httpVersion := "HTTP/1.1"
	if req != nil {
		httpVersion = req.httpVersion
	}
	return &Response{
		httpVersion: httpVersion,
		headers:     utils.NewKeyValueStore[string, string](),
	}
}

func (res *Response) BodyAsString() string {
	var bodySting string

	switch body := res.body.(type) {
	case string:
		bodySting = body

	case int:
		bodySting = strconv.Itoa(body)

	case float64:
		bodySting = strconv.FormatFloat(body, 'f', -1, 64)

	case float32:
		bodySting = strconv.FormatFloat(float64(body), 'f', -1, 32)

	case []byte:
		bodySting = string(body)

	case nil:
		bodySting = ""

	default:
		_bytes, err := json.Marshal(body)
		if err != nil {
			bodySting = fmt.Sprint("Error converting body to JSON -> Byte: ", err.Error())
		}
		bodySting = string(_bytes)
	}

	return bodySting
}

func (res *Response) SetDefaultHeaders(bodyStr string) {
	res.SetHeader("Content-Type", "application/json")
	res.SetHeader("Content-Length", strconv.Itoa(len(bodyStr)))
}

func (res *Response) HeadersToString() string {
	var headerString string
	for key, value := range res.headers.GetAll() {
		headerString += fmt.Sprintf("%v:%v\r\n", key, value)
	}
	headerString += "\r\n"

	return headerString
}

func (res *Response) SetStatus(code int) IResponse {

	// if no code is set, the 200 code is passed by default.
	if code == 0 {
		code = status.HTTP_200_OK
	}

	res.statusCode = code
	return res
}

func (res *Response) SetHeader(key string, value string) IResponse {
	res.headers.Set(key, value)
	return res
}

func (res *Response) Headers() *utils.KeyValueStore[string, string] {
	return res.headers
}

func (res *Response) HTTPVersion() string {
	return res.httpVersion
}

func (res *Response) StatusCode() int {
	return res.statusCode
}

func (res *Response) Body() any {
	return res.body
}

func (res *Response) Send(body any) IResponse {
	res.body = body
	return res
}

func (res *Response) GetResponseByte(isHead bool) []byte {
	bodyStr := res.BodyAsString()
	res.SetDefaultHeaders(bodyStr)
	responseString := res.httpVersion + " " + status.GetStatusString(res.statusCode) + "\r\n" + res.HeadersToString()

	if isHead {
		return []byte(responseString)
	}
	return []byte(responseString + bodyStr)
}
