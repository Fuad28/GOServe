package main

import (
	"github.com/codecrafters-io/http-server-starter-go/app/goserve"
	"github.com/codecrafters-io/http-server-starter-go/app/status"
)

// Test middlewares:
func middleware1(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	res.SetHeader("middleware1", "true")
	return req.Next(res)
}
func middleware2(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	res.SetHeader("middleware2", "true")
	res = req.Next(res)
	// We can access the final response in a middleware e.g fmt.Println(res.GetBody())
	return res
}

// Test routes
func tasks(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"handler": "tasks",
			"header":  res.GetHeaders().GetAll(),
		},
	)
}
func taskDetails(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"handler":    "taskDetails",
			"pathParams": req.PathParams.GetAll(),
		},
	)
}
func createTask(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"handler": "createTask",
			"body":    req.Body,
			"header":  res.GetHeaders().GetAll(),
		},
	)
}
func updateTask(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"handler":    "updateTask",
			"pathParams": req.PathParams.GetAll(),
			"body":       req.Body,
			"qParams":    req.QueryParams.GetAll(),
		},
	)
}
func deleteTask(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"handler":    "deleteTask",
			"pathParams": req.PathParams,
		},
	)
}

func main() {
	config := goserve.Config{
		Port:           4221,
		MaxRequestSize: goserve.ONE_MB,
	}

	route := goserve.NewServer(config)
	route.AddMiddleWare(goserve.CORSMiddleware(route.AllowedOrigins))
	route.AddMiddleWare(middleware1)

	route.GET("/tasks", tasks, middleware2)
	route.HEAD("/tasks", tasks)
	route.POST("/tasks", createTask)
	route.GET("/tasks/:id/details", taskDetails)
	route.PATCH("/tasks/:id/details", updateTask)
	route.DELETE("/tasks/:id", deleteTask)

	route.ServeAndListen()
}
