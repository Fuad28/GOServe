package main

import (
	"github.com/Fuad28/GOServe.git/goserve"
	"github.com/Fuad28/GOServe.git/goserve/status"
)

// Middlewares can either be added to the main router as in loggingMiddleware or to a single route as in the case of authenticationMiddlware
// Middlewares have to pass control to the next handler via req.Next(res)
// We can access the request and response at any point in the request lifecycle as in loggingMiddleware
// We can modify
// Multiple middlerwares can be passed at once.

// We can access the final response in a middleware.
func loggingMiddleware(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	// request logging logic

	res = req.Next(res)

	// response logging logic

	return res
}

func authenticationMiddlware(req *goserve.Request, res goserve.IResponse) goserve.IResponse {

	if token, exists := req.Headers().Get("Authorization"); exists {

		// Token authentication logic
		userId := token
		req.Store.Set("userId", userId)

	} else {
		return res.SetStatus(status.HTTP_401_UNAUTHORIZED).Send(
			goserve.JSON{"message": "unauthorized"},
		)
	}

	// pass control to the next handler in the handlerChain
	return req.Next(res)
}

func cacheMiddlware(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	res = req.Next(res)

	// response caching lo logic

	return res
}

func getTasksHandler(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"message": "tasks",
		},
	)
}

func main() {

	server := goserve.NewServer(goserve.Config{
		Port: 8000,
	})

	server.AddMiddleWare(loggingMiddleware)

	server.GET("/tasks", getTasksHandler, authenticationMiddlware, cacheMiddlware)

	server.ServeAndListen()
}
