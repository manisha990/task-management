package main

import (
    "github.com/gin-gonic/gin"    // Import Gin framework
    "database/sql"                // Import SQL package
	"strconv"
    _ "github.com/mattn/go-sqlite3" // Import SQLite3 driver
    "net/http"                      // Import HTTP package
    "time"                          // Import Time package
)

// Task represents the structure of a task
type Task struct {
    ID          int       `json:"id"`          // Unique identifier for each task
    Title       string    `json:"title"`       // Title or name of the task
    Description string    `json:"description"` // Brief description of the task
    DueDate     time.Time `json:"due_date"`    // Deadline or due date for the task
    Status      string    `json:"status"`      // Status of the task (e.g., pending, completed, in progress)
}

var db *sql.DB // Database connection variable

func main() {
    // Connect to SQLite database
    var err error
    db, err = sql.Open("sqlite3", "./tasks.db") // Open SQLite database
    if err != nil {
        panic(err)
    }
    defer db.Close() // Close database connection when main function exits

    // Create tasks table if it doesn't exist
    createTable()

    // Initialize Gin router
    router := gin.Default() // Create a new Gin router with default middleware

    // Define API routes
    router.POST("/tasks", createTask)   // Endpoint for creating a new task
    router.GET("/tasks/:id", getTaskByID) // Endpoint for retrieving a task by ID
    router.PUT("/tasks/:id", updateTask) // Endpoint for updating a task by ID
    router.DELETE("/tasks/:id", deleteTask) // Endpoint for deleting a task by ID
    router.GET("/tasks", getAllTasks)   // Endpoint for retrieving all tasks

    // Start server
    router.Run(":8080") // Start the HTTP server and listen on port 8080
}

// createTable creates tasks table if it doesn't exist
func createTable() {
    createTableSQL := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        description TEXT,
        due_date TIMESTAMP,
        status TEXT
    );` // SQL query to create tasks table with specified schema
    _, err := db.Exec(createTableSQL) // Execute SQL query
    if err != nil {
        panic(err) // Panic if there is an error
    }
}

// createTask handles POST request to create a new task
func createTask(c *gin.Context) {
    var task Task // Create a new Task variable
    if err := c.ShouldBindJSON(&task); err != nil { // Bind JSON payload to Task struct
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    insertTaskSQL := "INSERT INTO tasks (title, description, due_date, status) VALUES (?, ?, ?, ?)"
    result, err := db.Exec(insertTaskSQL, task.Title, task.Description, task.DueDate, task.Status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    taskID, _ := result.LastInsertId() // Get last inserted ID
    task.ID = int(taskID)              // Convert ID to integer
    c.JSON(http.StatusCreated, task)   // Return created task with status code 201 (Created)
}

// getTaskByID handles GET request to retrieve a task by ID
func getTaskByID(c *gin.Context) {
    id := c.Param("id") // Get task ID from request parameter

    var task Task // Create a new Task variable
    row := db.QueryRow("SELECT id, title, description, due_date, status FROM tasks WHERE id = ?", id)
    err := row.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"}) // Return error if task not found
        return
    }

    c.JSON(http.StatusOK, task) // Return task details with status code 200 (OK)
}

// updateTask handles PUT request to update a task by ID
func updateTask(c *gin.Context) {
    id := c.Param("id") // Get task ID from request parameter

    var task Task // Create a new Task variable
    if err := c.ShouldBindJSON(&task); err != nil { // Bind JSON payload to Task struct
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updateTaskSQL := "UPDATE tasks SET title=?, description=?, due_date=?, status=? WHERE id=?"
    _, err := db.Exec(updateTaskSQL, task.Title, task.Description, task.DueDate, task.Status, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    task.ID, _ = strconv.Atoi(id) // Convert ID to integer
    c.JSON(http.StatusOK, task)   // Return updated task with status code 200 (OK)
}

// deleteTask handles DELETE request to delete a task by ID
func deleteTask(c *gin.Context) {
    id := c.Param("id") // Get task ID from request parameter

    deleteTaskSQL := "DELETE FROM tasks WHERE id=?"
    _, err := db.Exec(deleteTaskSQL, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"}) // Return success message with status code 200 (OK)
}

// getAllTasks handles GET request to retrieve all tasks
func getAllTasks(c *gin.Context) {
    rows, err := db.Query("SELECT id, title, description, due_date, status FROM tasks")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var tasks []Task // Create a slice to store tasks
    for rows.Next() {
        var task Task // Create a new Task variable
        err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        tasks = append(tasks, task) // Append task to tasks slice
    }

    c.JSON(http.StatusOK, tasks) // Return list of tasks with status code 200 (OK)

}