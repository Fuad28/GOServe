package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

type Request struct {
	HTTPVersion string
	Headers     map[string]string
	Body        interface{}
	Method      string
	Target      string
}

func NewRequest(req string) (*Request, error) {
	scanner := bufio.NewScanner(strings.NewReader(req))

	if !scanner.Scan() {
		return nil, errors.New("invalid request: missing request line")
	}

	request := Request{}

	// Parse request line
	requestLine := strings.Fields(scanner.Text())
	request.Method = requestLine[0]
	request.Target = requestLine[1]
	request.HTTPVersion = requestLine[2]

	// Parse headers
	headers := map[string]string{}
	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			break
		}

		headerParts := strings.Split(text, ": ")

		if len(headerParts) != 2 {
			return nil, errors.New("invalid request: invalid header")

		} else {
			headers[headerParts[0]] = headerParts[1]
		}

	}
	request.Headers = headers

	// Parse body
	if scanner.Scan() {
		request.Body = scanner.Text()
	}
	return &request, nil
}

type Response struct {
	HTTPVersion string
	StatusCode  int
	Headers     map[string]any
	Body        string
}

func (res *Response) GetResponseByte() ([]byte, error) {
	var response string
	CRLF := "\r\n"

	if res.StatusCode == 0 {
		return []byte{}, errors.New("status code is required")
	}

	status, err := getStatus(strconv.Itoa(res.StatusCode))

	if err != nil {
		return []byte{}, err
	}

	http_version := res.HTTPVersion
	if http_version == "" {
		http_version = "HTTP/1.1"
	}

	response = http_version + " " + status + CRLF

	if res.Body != "" {
		res.Headers = map[string]any{"Content-Type": "text/plain", "Content-Length": len(res.Body)}
		response += "Content-Type:text/plain" + CRLF + "Content-Length:" + " " + strconv.Itoa(len(res.Body)) + CRLF + CRLF + res.Body

	} else {
		response += CRLF
	}
	fmt.Println(response)
	return []byte(response), nil
}

type HTTPStatus struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	Descriptions string `json:"description"`
}

type HTTPStatusMap map[string]HTTPStatus

func getStatus(code string) (string, error) {
	httpStatues := HTTPStatusMap{}
	err := utils.LoadFile[HTTPStatusMap]("http_statuses.json", &httpStatues)

	if err != nil {
		return "", err
	}
	status := httpStatues[code]
	statusString := strconv.Itoa(status.Code) + " " + status.Message
	return statusString, nil

}

func handleError(conn *net.Conn, err error, client bool) {
	c := *conn
	if err != nil {
		defer c.Close()
		defer os.Exit(1)

		if client {
			c.Write([]byte(err.Error()))

		} else {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		handleError(&conn, err, false)

		req := make([]byte, 1024)
		conn.Read(req)

		request, err := NewRequest(string(req))
		handleError(&conn, err, true)

		wg := sync.WaitGroup{}
		wg.Add(1)
		defer wg.Wait()

		go func(conn *net.Conn) {
			defer wg.Done()

			res, err := handleRequest(request)
			handleError(conn, err, true)

			responseByte, err := res.GetResponseByte()
			handleError(conn, err, false)

			(*conn).Write(responseByte)
			(*conn).Close()
		}(&conn)
	}
}

func handleRequest(req *Request) (Response, error) {
	response := Response{}
	if strings.HasPrefix(req.Target, "/echo") {
		response.StatusCode = 200
		response.Body = strings.SplitN(req.Target, "/echo/", 2)[1]

	} else if strings.HasPrefix(req.Target, "/user-agent") {
		response.StatusCode = 200
		response.Body = req.Headers["User-Agent"]

	} else if req.Target == "/" {
		response.StatusCode = 200

	} else {
		response.StatusCode = 404

	}
	return response, nil

}
