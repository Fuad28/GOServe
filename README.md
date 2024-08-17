# GoServe

**GOServe** is a lightweight and extensible HTTP server framework written in Go, designed to simplify the process of building JSON REST APIs. It’s ideal for projects requiring a minimalistic and efficient server.

## Table of Contents

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Configuration](#configuration)
5. [Routing](#routing)
6. [Middleware](#middleware)
7. [CORS Support](#cors-support)
8. [Passing Data Around](#passing-data-around)
9. [Contributing](#contributing)
10. [License](#license)
11. [Contributors](#contributors)


## Features

- **Lightweight**: Minimal dependencies, designed to be fast and efficient.
- **JSON-Native**: Built-in support for JSON request and response bodies.
- **Flexible Routing**: Supports dynamic routing and path parameters.
- **Middleware Support**: Easily extend functionality with custom middleware.
- **CORS Handling**: Built-in support for Cross-Origin Resource Sharing (CORS).
- **Error Handling**: Simplified error handling with customizable responses.

## Installation

To install GOServe, use `go get`:

```bash
go get github.com/fuad28/goserve
```


### Quick Start
Here's a simple example to get you started with GOServe:

```go
package main

import (
    "github.com/fuad28/goserve"
    "github.com/fuad28/goserve/status"
)

func main() {
    server := goserve.NewServer(goserve.Config{
        Port: 8000,
    })

    server.Get("/hello", func(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
        return res.SetStatus(status.HTTP_200_OK).Send("Hello World!")
    })

    server.ListenAndServe()
}
```
**Full examples can be found in the example folder.**

### Configuration
GOServe can be configured using the `goserve.Config` struct. Key options include:

- **Port**: The port on which the server listens defaults to 8000.
- **MaxRequestSize**: Maximum size of the request body defaults to 1MB.
- **AllowedOrigins**: Origins allowed for CORS.

Example:

```go
server := goserve.NewServer(goserve.Config{
    Port:           8000,
    MaxRequestSize: goserve.ONE_MB * 10,
    AllowedOrigins: []string{"http://localhost:3000"},
})
```


### Routing
Define routes with `Get`, `Post`, `Put`, `Delete`, and other HTTP methods. Routes support dynamic path and query parameters.

Example:

```go
server.Get("/tasks/:id", func(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
    userID := req.PathParams().Get("id")
    queryParameters := req.QueryParams().GetAll()
    // Logic to fetch task by ID
    return res.SetStatus(status.HTTP_200_OK).Send(goserve.JSON{"task": task})
})
```


### Middleware
Middleware allows you to extend functionality with custom middleware easily. In a middleware, you have access to the request and response throughout the request-response lifecycle. Use middleware to implement logging, authentication, etc.
* **Note: Ensure to call the req.Next() method in your middleware function to pass control to the next handler in the handlerChain**

Example:

```go
func loggingMiddleware(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	// request logging logic

	res = req.Next(res)

	// response logging logic

	return res
}

server.AddMiddlewares(loggingMiddleware)
```


### CORS Support
GOServe has built-in CORS support, configurable via middleware. Allow specific origins, methods, and headers.
* **Note: When testing with non-browser clients like Postman and Curl, you have to set your Origin header manually.

Example:

```go
server.AddMiddlewares(goserve.CORSMiddleware([]string{"http://example.com"}))
```

### Passing data around
GOServe requests have Store attribute that allows you to pass and retrieve data throughout the request lifecycle

Example:

```go
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

server.AddMiddlewares(authenticationMiddlware)

```

In your handler, you would access the userId as:

```go
func allTasksHandler(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
  userId, exists := req.Store.Get("userId", userId)
}
```


### Contributing
Contributions are welcome! Please read the [contributing guide](./contributing.md) to learn about our development process, how to propose bug fixes and improvements, and how to build and test your changes to GOServe.

Check out the currently open [issues](https://github.com/fuad28/goserve/issues).

### License
This project is licensed under the [MIT license](./LICENSE.md).

### Contributors
<a href="https://github.com/fuad28/goserve/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=fuad28/goserve" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

