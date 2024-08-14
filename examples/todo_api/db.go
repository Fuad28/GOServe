package main

import "github.com/Fuad28/GOServe.git/goserve/utils"

// Data setup
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	UserId int    `json:"userId"`
}

// Initialize our "database" by adding some data to the "tables"
func initDB(users *utils.KeyValueStore[int, User], tasks *utils.KeyValueStore[int, Task]) {
	users.Set(1, User{Id: 1, Name: "Alice"})
	users.Set(2, User{Id: 2, Name: "Bob"})

	tasks.Set(1, Task{Id: 1, Title: "Buy groceries", UserId: 1})
	tasks.Set(2, Task{Id: 2, Title: "Read a book", UserId: 2})
}
