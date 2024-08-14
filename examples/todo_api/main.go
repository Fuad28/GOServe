package main

import (
	"github.com/Fuad28/GOServe.git/goserve"
	"github.com/Fuad28/GOServe.git/goserve/utils"
)

// This example builds a simple ToDo api to demonstrate GOServe usage.
// We use the utils.KeyValueStore provided by GOServe as database
// Since authenticationMiddlware sample, directly sets the token as the userId, you should set the Authorization header to 1 or 2.

var users *utils.KeyValueStore[int, User]
var tasks *utils.KeyValueStore[int, Task]

func main() {

	// Configure server
	server := goserve.NewServer(goserve.Config{
		Port:           8000,
		MaxRequestSize: goserve.ONE_MB * 10,
		AllowedOrigins: []string{"http://127.0.0.1"},
	})

	// initialze database with dummy data
	initDB(users, tasks)

	// Add server level middlewares
	server.AddMiddleWare(goserve.CORSMiddleware(server.AllowedOrigins()))
	server.AddMiddleWare(authenticationMiddlware)

	// Register routes
	server.GET("/tasks", allTasks)
	server.POST("/tasks", createTask)
	server.GET("/tasks/:id", taskDetails)
	server.DELETE("/tasks/:id", deleteTask)

	// Start server and listen for connections
	server.ServeAndListen()
}
