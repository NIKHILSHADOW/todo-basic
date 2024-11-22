package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct {
	TaskId    string `json:"id"`
	TaskName  string `json:"title"`
	Completed bool   `json:"status"`
}

var todos = []todo{
	{TaskId: "1", TaskName: "read books", Completed: false},
	{TaskId: "2", TaskName: "sleep", Completed: false},
	{TaskId: "3", TaskName: "code", Completed: true},
	{TaskId: "4", TaskName: "eat", Completed: true},
}

func getTask(id string) (*todo, error) {

	for i, task := range todos {

		if task.TaskId == id {
			return &todos[i], nil
		}
	}

	return nil, errors.New("task not found")
}

func getTaskById(context *gin.Context) {

	id := context.Param("id")
	task, err := getTask(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, errors.New("Task with id "+id+" not found"))
		return
	}

	context.IndentedJSON(http.StatusOK, task)

}

func updateTask(context *gin.Context) {

	var taskReq todo
	id := context.Param("id")
	if err := context.Bind(&taskReq); err != nil {
		context.IndentedJSON(http.StatusNotFound, errors.New("given task is not valid"))
		return
	}

	task, err := getTask(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, err)
		return
	}

	if taskReq.TaskName != "" {
		task.TaskName = taskReq.TaskName
	}

	if taskReq.Completed != task.Completed {
		task.Completed = taskReq.Completed
	}

	context.IndentedJSON(http.StatusOK, task)
}

func deleteTask(context *gin.Context) {

	id := context.Param("id")
	j := 0

	for _, task := range todos {

		if task.TaskId != id {

			todos[j].TaskId = task.TaskId
			todos[j].TaskName = task.TaskName
			todos[j].Completed = task.Completed

			j++
		}
	}

	todos = todos[:j]

	context.IndentedJSON(http.StatusOK, "deleted")

}

func main() {
	router := gin.Default()
	router.GET("/todos", getTasks)
	router.POST("/todos", addTask)
	router.GET("/todos/:id", getTaskById)
	router.PATCH("/todos/:id", updateTask)
	router.DELETE("/todos/:id", deleteTask)
	router.Run("localhost:9090")

}

func addTask(context *gin.Context) {

	var task todo

	if err := context.Bind(&task); err != nil {
		return
	}

	todos = append(todos, task)

	context.IndentedJSON(http.StatusCreated, task)
}

func getTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}
