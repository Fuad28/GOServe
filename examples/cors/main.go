package main

import (
	"github.com/Fuad28/GOServe.git/goserve"
	"github.com/Fuad28/GOServe.git/goserve/status"
)

// GOServe provides a CORSMiddleware

func getTasksHandler(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"message": "tasks",
		},
	)
}

func main() {

	server := goserve.NewServer(goserve.Config{
		MaxRequestSize: goserve.ONE_MB,
		AllowedOrigins: []string{"https://www.google.com/"},
	})

	server.AddMiddleWare(goserve.CORSMiddleware(server.AllowedOrigins()))

	server.GET("/tasks", getTasksHandler)

	server.ServeAndListen()
}
