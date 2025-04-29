package main

import (
	"fmt"
	"task_manager/db"
	"task_manager/handlers"
	"time"

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

	loc, err := time.LoadLocation("Europe/Madrid") // zona horaria deseada, tío
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	fmt.Println("Hora local en Madrid:", now)
}
