package controllers

import (
	"net/http"
	"task-management/initializers"
	"task-management/models"

	"github.com/gin-gonic/gin"
)

func TaskCreate(c *gin.Context) {
	//Get data off req body
	var body struct {
		Title       string
		Description string
		DueDate     string
		Status      string
	}
	c.Bind(&body)

	//create a post
	task := models.Task{Title: body.Title, Description: body.Description, DueDate: body.DueDate, Status: body.Status}
	result := initializers.DB.Create(&task)

	if result.Error != nil {
		c.Status(400)
		return
	}

	//Return it
	c.JSON(200, gin.H{
		"task":    task,
		"message": "success",
	})
}

func TaskIndex(c *gin.Context) {
	//Get the tasks
	var tasks []models.Task
	initializers.DB.Find(&tasks)

	//respond with them

	c.JSON(200, gin.H{
		"tasks":   tasks,
		"message": "success",
	})

}

func GetTask(c *gin.Context) {

	// Get the task ID from the URL parameters
	taskID := c.Param("id")

	// Check if the task ID is provided
	//  if taskID == "" {
	// 	 c.JSON(http.StatusBadRequest, gin.H{
	// 		 "error":   "Task ID is required",
	// 		 "message": "failure",
	// 	 })
	// 	 return
	//  }

	// Query the database for the task with the specified ID
	var task models.Task
	result := initializers.DB.First(&task, taskID)

	// Check if the task is found
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Task not found",
			"message": "failure",
		})
		return
	}

	//respond with them
	c.JSON(200, gin.H{
		"tasks":   task,
		"message": "success",
	})

}

func TaskUpdate(c *gin.Context) {
	//get the id off the url
	id := c.Param("id")

	//get the data of the req body
	var body struct {
		Title       string
		Description string
		DueDate     string
		Status      string
	}
	c.Bind(&body)

	//find the task we are updating
	var task models.Task
	initializers.DB.First(&task, id)

	//update

	initializers.DB.Model(&task).Updates(models.Task{Title: body.Title, Description: body.Description, DueDate: body.DueDate, Status: body.Status})

	//respond

	c.JSON(200, gin.H{
		"tasks":   task,
		"message": "Updated successfully ",
	})

}
func TaskDelete(c *gin.Context) {
	//get the id of the url
	id := c.Param("id")

	//delete the post
	initializers.DB.Delete(&models.Task{}, id)

	//respond
	c.JSON(200, gin.H{

		"message": "deleted successfully ",
	})

}
