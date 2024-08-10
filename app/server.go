package main

import (
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/server"
	"github.com/codecrafters-io/http-server-starter-go/app/status"
	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

// Test middlewares:
func middleware1(req *server.Request, res server.IResponse) server.IResponse {
	res.SetHeader("middleware1", "true")
	return req.Next(res)
}
func middleware2(req *server.Request, res server.IResponse) server.IResponse {
	res.SetHeader("middleware2", "true")
	res = req.Next(res)
	// We can access the final response in a middleware e.g fmt.Println(res.GetBody())
	return res
}

// Test routes
func tasks(req *server.Request, res server.IResponse) server.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		server.JSON{
			"handler": "tasks",
			"header":  res.GetHeaders().GetAll(),
		},
	)
}
func taskDetails(req *server.Request, res server.IResponse) server.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		server.JSON{
			"handler":    "taskDetails",
			"pathParams": req.PathParams.GetAll(),
		},
	)
}
func createTask(req *server.Request, res server.IResponse) server.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		server.JSON{
			"handler": "createTask",
			"body":    req.Body,
			"header":  res.GetHeaders().GetAll(),
		},
	)
}
func updateTask(req *server.Request, res server.IResponse) server.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		server.JSON{
			"handler":    "updateTask",
			"pathParams": req.PathParams.GetAll(),
			"body":       req.Body,
			"qParams":    req.QueryParams.GetAll(),
		},
	)
}
func deleteTask(req *server.Request, res server.IResponse) server.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		server.JSON{
			"handler":    "deleteTask",
			"pathParams": req.PathParams,
		},
	)
}

// codecrafters routes
func home(req *server.Request, res server.IResponse) server.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(nil)
}
func echo(req *server.Request, res server.IResponse) server.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(strings.SplitN(req.Path, "/echo/", 2)[1])
}
func useragent(req *server.Request, res server.IResponse) server.IResponse {
	agent, _ := req.Headers.Get("User-Agent")
	return res.SetStatus(status.HTTP_200_OK).Send(agent)
}
func files(req *server.Request, res server.IResponse) server.IResponse {
	fileName := strings.SplitN(req.Path, "/files/", 2)[1]
	fileLocation := "/tmp/" + fileName
	fileContent, err := utils.LoadPlainTextFile(fileLocation)
	statusCode := 200

	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			statusCode = 404
		} else {
			statusCode = 400
		}

	}
	return res.SetStatus(statusCode).SetHeader("Content-Type", "application/octet-stream").Send(fileContent)
}

func main() {
	config := server.Config{
		Port:           4221,
		MaxRequestSize: server.ONE_MB,
	}
	route := server.NewServer(config)
	// route.AddMiddleWare(server.CORSMiddleware([]string{"0.0.0.1"}))
	route.AddMiddleWare(middleware1)

	route.GET("/tasks", tasks, middleware2)
	route.HEAD("/tasks", tasks)
	route.POST("/tasks", createTask)
	route.GET("/tasks/:id/details", taskDetails)
	route.PATCH("/tasks/:id/details", updateTask)
	route.DELETE("/tasks/:id", deleteTask)

	// codecrafters controllers
	route.GET("/", home)
	route.GET("/echo", echo)
	route.GET("/user-agent", useragent)
	route.GET("/file", files)

	route.ServeAndListen()
}

// import (
// 	"bufio"
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"net"
// 	"os"
// 	"strconv"
// 	"strings"

// 	"github.com/codecrafters-io/http-server-starter-go/app/utils"
// 	// "github.com/codecrafters-io/http-server-starter-go/app/server"
// )

// // Sample usage
// type Request struct {
// 	HTTPVersion string
// 	Headers     map[string]string
// 	Body        interface{}
// 	Method      string
// 	Target      string
// }

// func NewRequest(req string) (*Request, error) {
// 	scanner := bufio.NewScanner(strings.NewReader(req))

// 	if !scanner.Scan() {
// 		return nil, errors.New("invalid request: missing request line")
// 	}

// 	request := Request{}

// 	// Parse request line
// 	requestLine := strings.Fields(scanner.Text())
// 	request.Method = requestLine[0]
// 	request.Target = requestLine[1]
// 	request.HTTPVersion = requestLine[2]

// 	// Parse headers
// 	headers := map[string]string{}
// 	for scanner.Scan() {
// 		text := scanner.Text()

// 		if text == "" {
// 			break
// 		}

// 		headerParts := strings.Split(text, ": ")

// 		if len(headerParts) != 2 {
// 			return nil, errors.New("invalid request: invalid header")

// 		} else {
// 			headers[headerParts[0]] = headerParts[1]
// 		}

// 	}
// 	request.Headers = headers

// 	// Parse body
// 	if scanner.Scan() {
// 		request.Body = scanner.Text()
// 	}
// 	return &request, nil
// }

// type Response struct {
// 	HTTPVersion   string
// 	StatusCode    int
// 	Headers       map[string]any
// 	HeadersString string
// 	Body          string
// 	BodyType      string
// }

// func (res *Response) GetResponseByte() ([]byte, error) {

// 	CRLF := "\r\n"

// 	res.SetHeaders()

// 	var response string

// 	if res.StatusCode == 0 {
// 		return []byte{}, errors.New("status code is required")
// 	}

// 	status, err := getStatus(strconv.Itoa(res.StatusCode))

// 	if err != nil {
// 		return []byte{}, err
// 	}

// 	response = res.HTTPVersion + " " + status + CRLF + res.HeadersString + res.Body

// 	return []byte(response), nil
// }

// func (res *Response) SetHeaders() {
// 	CRLF := "\r\n"

// 	if res.BodyType == "" {
// 		res.BodyType = "text/plain"
// 	}

// 	if res.Body != "" {
// 		contentLength := strconv.Itoa(len(res.Body))
// 		res.Headers = map[string]any{"Content-Type": res.BodyType, "Content-Length": len(res.Body)}
// 		res.HeadersString = fmt.Sprintf("Content-Type:%v"+CRLF+"Content-Length: %v"+CRLF+CRLF, res.BodyType, contentLength)

// 	} else {
// 		res.HeadersString += CRLF
// 	}
// }

// type HTTPStatus struct {
// 	Code         int    `json:"code"`
// 	Message      string `json:"message"`
// 	Descriptions string `json:"description"`
// }

// type HTTPStatusMap map[string]HTTPStatus

// func getStatus(code string) (string, error) {
// 	httpStatues := HTTPStatusMap{}
// 	err := utils.LoadFile[HTTPStatusMap]("http_statuses.json", &httpStatues)

// 	if err != nil {
// 		return "", err
// 	}
// 	status := httpStatues[code]
// 	statusString := strconv.Itoa(status.Code) + " " + status.Message
// 	return statusString, nil

// }

// func handleError(conn *net.Conn, err error, client bool) {
// 	c := *conn
// 	if err != nil {
// 		defer c.Close()

// 		if client {
// 			c.Write([]byte(err.Error()))

// 		} else {
// 			fmt.Println(err.Error())
// 		}
// 	}
// }

// func main() {
// 	fmt.Println("Logs from your program will appear here!")

// 	tempFileDirectory := flag.String("directory", "/tmp/", "directory to find files.")
// 	flag.Parse()

// 	l, err := net.Listen("tcp", "0.0.0.0:4221")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 4221: ", err.Error())
// 		os.Exit(1)
// 	}

// 	for {
// 		conn, err := l.Accept()
// 		handleError(&conn, err, true)

// 		go func() {
// 			req := make([]byte, 1024)
// 			_, err = conn.Read(req)
// 			handleError(&conn, err, true)

// 			request, err := NewRequest(string(req))
// 			handleError(&conn, err, true)

// 			res, err := handleRequest(request, tempFileDirectory)
// 			handleError(&conn, err, true)

// 			responseByte, err := res.GetResponseByte()
// 			handleError(&conn, err, false)

// 			conn.Write(responseByte)
// 			conn.Close()
// 		}()
// 	}
// }

// func handleRequest(req *Request, tempFileDirectory *string) (Response, error) {
// 	response := Response{
// 		HTTPVersion: req.HTTPVersion,
// 	}

// 	if strings.HasPrefix(req.Target, "/echo") {
// 		response.StatusCode = 200
// 		response.Body = strings.SplitN(req.Target, "/echo/", 2)[1]

// 	} else if strings.HasPrefix(req.Target, "/files/") {

// 		fileName := strings.SplitN(req.Target, "/files/", 2)[1]
// 		fileLocation := *tempFileDirectory + fileName

// 		fileContent, err := utils.LoadPlainTextFile(fileLocation)

// 		if err != nil {
// 			if _, ok := err.(*os.PathError); ok {
// 				response.StatusCode = 404
// 			} else {
// 				response.StatusCode = 400
// 			}

// 		} else {
// 			response.StatusCode = 200
// 			response.Body = fileContent
// 			response.BodyType = "application/octet-stream"
// 		}

// 	} else if strings.HasPrefix(req.Target, "/user-agent") {
// 		response.StatusCode = 200
// 		response.Body = req.Headers["User-Agent"]

// 	} else if req.Target == "/" {
// 		response.StatusCode = 200

// 	} else {
// 		response.StatusCode = 404

// 	}
// 	return response, nil

// }
