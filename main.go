package main

import (
	"task-management/controllers"
	"task-management/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {
	// routes 
	r := gin.Default()
	//create task 
	r.POST("/create-task", controllers.TaskCreate)

	//retrive tasks
	r.GET("/get-tasks", controllers.TaskIndex)

	//update task
	r.PUT("/update-task/:id", controllers.TaskUpdate)

	//get specific task
	r.GET("get-task/:id", controllers.GetTask)
	
	//delete task
	r.DELETE("/task/:id", controllers.TaskDelete)
	r.Run("localhost:3000") // listen and serve on 0.0.0.0:8080
}
