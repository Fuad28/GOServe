package goserve

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Fuad28/GOServe.git/goserve/status"
	"github.com/Fuad28/GOServe.git/goserve/utils"
)

type IResponse interface {
	SetStatus(int) IResponse
	SetHeader(key string, value string) IResponse
	Headers() *utils.KeyValueStore[string, string]
	Body() any
	Send(any) IResponse
	GetResponseByte(bool) []byte
}

type Response struct {
	HTTPVersion string
	StatusCode  int
	headers     *utils.KeyValueStore[string, string]
	body        any
}

func NewResponse(req *Request) *Response {
	httpVersion := "HTTP/1.1"
	if req != nil {
		httpVersion = req.httpVersion
	}
	return &Response{
		HTTPVersion: httpVersion,
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

	res.StatusCode = code
	return res
}

func (res *Response) SetHeader(key string, value string) IResponse {
	res.headers.Set(key, value)
	return res
}

func (res *Response) Headers() *utils.KeyValueStore[string, string] {
	return res.headers
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
	responseString := res.HTTPVersion + " " + status.GetStatusString(res.StatusCode) + "\r\n" + res.HeadersToString()

	if isHead {
		return []byte(responseString)
	}
	return []byte(responseString + bodyStr)
}
