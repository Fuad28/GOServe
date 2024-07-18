package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

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
	HTTPVersion   string
	StatusCode    int
	Headers       map[string]any
	HeadersString string
	Body          string
	BodyType      string
}

func (res *Response) GetResponseByte() ([]byte, error) {

	CRLF := "\r\n"

	res.SetHeaders()

	var response string

	if res.StatusCode == 0 {
		return []byte{}, errors.New("status code is required")
	}

	status, err := getStatus(strconv.Itoa(res.StatusCode))

	if err != nil {
		return []byte{}, err
	}

	response = res.HTTPVersion + " " + status + CRLF + res.HeadersString + res.Body

	return []byte(response), nil
}

func (res *Response) SetHeaders() {
	CRLF := "\r\n"

	if res.BodyType == "" {
		res.BodyType = "text/plain"
	}

	if res.Body != "" {
		contentLength := strconv.Itoa(len(res.Body))
		res.Headers = map[string]any{"Content-Type": res.BodyType, "Content-Length": len(res.Body)}
		res.HeadersString = fmt.Sprintf("Content-Type:%v"+CRLF+"Content-Length: %v"+CRLF+CRLF, res.BodyType, contentLength)

	} else {
		res.HeadersString += CRLF
	}
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

		if client {
			c.Write([]byte(err.Error()))

		} else {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	tempFileDirectory := flag.String("directory", "/tmp/", "directory to find files.")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221: ", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		handleError(&conn, err, true)

		go func() {
			req := make([]byte, 1024)
			_, err = conn.Read(req)
			handleError(&conn, err, true)

			request, err := NewRequest(string(req))
			handleError(&conn, err, true)

			res, err := handleRequest(request, tempFileDirectory)
			handleError(&conn, err, true)

			responseByte, err := res.GetResponseByte()
			handleError(&conn, err, false)

			conn.Write(responseByte)
			conn.Close()
		}()
	}
}

func handleRequest(req *Request, tempFileDirectory *string) (Response, error) {
	response := Response{
		HTTPVersion: req.HTTPVersion,
	}

	if strings.HasPrefix(req.Target, "/echo") {
		response.StatusCode = 200
		response.Body = strings.SplitN(req.Target, "/echo/", 2)[1]

	} else if strings.HasPrefix(req.Target, "/files/") {

		fileName := strings.SplitN(req.Target, "/files/", 2)[1]
		fileLocation := *tempFileDirectory + fileName

		fileContent, err := utils.LoadPlainTextFile(fileLocation)

		if err != nil {
			if _, ok := err.(*os.PathError); ok {
				fmt.Println("HEREEE")
				response.StatusCode = 404
			} else {
				response.StatusCode = 400
			}

		} else {
			response.StatusCode = 200
			response.Body = fileContent
			response.BodyType = "application/octet-stream"
		}

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
