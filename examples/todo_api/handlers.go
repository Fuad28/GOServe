package main

import (
	"strconv"

	"github.com/Fuad28/GOServe.git/goserve"
	"github.com/Fuad28/GOServe.git/goserve/status"
)

func allTasks(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"tasks": tasks.GetAll(),
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

	if task, exists := tasks.Get(taskId); exists && task.UserId == userId {
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
	// body := req.Body()
	// var task Task

	// err := json.Unmarshal([]byte(body), &task)
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"handler": "createTask",
			"body":    req.Body,
			"header":  res.Headers().GetAll(),
		},
	)
}

func updateTask(req *goserve.Request, res goserve.IResponse) goserve.IResponse {
	return res.SetStatus(status.HTTP_200_OK).Send(
		goserve.JSON{
			"handler":    "updateTask",
			"pathParams": req.PathParams().GetAll(),
			"body":       req.Body,
			"qParams":    req.QueryParams().GetAll(),
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
