package main

import (
	"github.com/Fuad28/GOServe.git/goserve"
	"github.com/Fuad28/GOServe.git/goserve/status"
)

func getTasksHandler(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send("Hello world")
}
func main() {
	server := goserve.NewServer(goserve.Config{
		Port: 8000,
	})

	server.GET("/tasks", getTasksHandler)

	server.ServeAndListen()
}
