package main

import (
	"task_manager/db"
	"task_manager/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Importa el driver de MySQL
)

func main() {
	db.Connect()

	r := gin.Default()

	r.POST("/api/tasks", handlers.CreateTask)
	r.GET("/api/tasks", handlers.GetTasks)
	r.GET("/api/tasks/:id", handlers.GetTasks)
	r.PUT("/api/tasks/:id", handlers.UpdateTask)
	r.DELETE("/api/tasks/:id", handlers.DeleteTask)

	r.Run(":8080") // Levanta el servidor en localhost:8080
}
