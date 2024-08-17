package main

import (
	"github.com/Fuad28/GOServe.git/goserve/utils"
)

// Data setup
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	Id     int    `json:"id"`
	Title  string `json:"title" valid:"required"`
	UserId int    `json:"userId"`
}

// Initialize our "database" by seeding some data to the "tables"
func seedDB(users *utils.KeyValueStore[int, User], tasks *utils.KeyValueStore[int, Task]) {
	users.Set(1, User{Id: 1, Name: "John"})
	users.Set(2, User{Id: 2, Name: "Doe"})

	tasks.Set(1, Task{Id: 1, Title: "Add subroute support", UserId: 1})
	tasks.Set(2, Task{Id: 2, Title: "Add HTTP encoding", UserId: 2})
}

func getTasksByUserId(tasks *utils.KeyValueStore[int, Task], userId int) *utils.KeyValueStore[int, Task] {
	userTasks := utils.NewKeyValueStore[int, Task]()

	for _, task := range tasks.GetAll() {
		if task.UserId == userId {
			userTasks.Set(task.Id, task)
		}
	}

	return userTasks
}
