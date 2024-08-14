package main

import (
	"fmt"
	"strconv"

	"github.com/Fuad28/GOServe.git/goserve"
	"github.com/Fuad28/GOServe.git/goserve/status"
)

func allTasks(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	// userId exists because the route is protected on server level
	userId, _ := req.Store.Get("userId")
	userTasks := getTasksByUserId(tasks, userId.(int))

	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"tasks": userTasks.GetAll(),
		},
	)
}

func taskDetails(req *goserve.Request, res goserve.IResponse) goserve.IResponse {

	// userId exists because the route is protected on server level
	userId, _ := req.Store.Get("userId")

	taskIdStr, _ := req.PathParams().Get("id")
	taskId, err := strconv.Atoi(taskIdStr)

	if err != nil {
		return res.SetStatus(status.HTTP_400_BAD_REQUEST).Send(
			goserve.JSON{
				"error": "Invalid id",
			},
		)
	}

	userTasks := getTasksByUserId(tasks, userId.(int))
	if task, exists := userTasks.Get(taskId); exists {
		return res.SetStatus(status.HTTP_200_OK).Send(
			goserve.JSON{
				"task": task,
			},
		)

	} else {
		return res.SetStatus(status.HTTP_404_NOT_FOUND).Send(
			goserve.JSON{
				"error": "Not Found",
			},
		)
	}

}

func createTask(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	var task Task

	if err := req.Body(&task); err != nil {
		return res.SetStatus(status.HTTP_400_BAD_REQUEST).Send(
			goserve.JSON{
				"error": fmt.Sprintf("invalid body: %v", err.Error()),
			},
		)
	}

	userId, _ := req.Store.Get("userId")
	task.Id = len(tasks.GetAll()) + 1
	task.UserId = userId.(int)

	return res.SetStatus(status.HTTP_201_CREATED).Send(
		goserve.JSON{
			"task": task,
		},
	)
}

func deleteTask(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	// userId exists because the route is protected on server level
	userId, _ := req.Store.Get("userId")

	taskIdStr, _ := req.PathParams().Get("id")
	taskId, err := strconv.Atoi(taskIdStr)

	if err != nil {
		return res.SetStatus(status.HTTP_400_BAD_REQUEST).Send(
			goserve.JSON{
				"error": "Invalid id",
			},
		)
	}

	userTasks := getTasksByUserId(tasks, userId.(int))
	if task, exists := userTasks.Get(taskId); exists && task.UserId == userId {
		tasks.Delete(taskId)
		return res.SetStatus(status.HTTP_204_NO_CONTENT).Send(nil)

	} else {
		return res.SetStatus(status.HTTP_404_NOT_FOUND).Send(
			goserve.JSON{
				"error": "Not Found",
			},
		)
	}

}
