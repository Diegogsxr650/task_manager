package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"task_manager/db"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task.Title == "" || task.Description == "" || task.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, description, and status are required"})
		return
	}

	// Validación de que el estado sea uno de los permitidos
	validStatuses := []string{"new", "ongoing", "completed"}
	isValidStatus := false
	for _, status := range validStatuses {
		if task.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Allowed values are 'pendiente', 'en progreso', or 'completado'"})
		return
	}

	_, err := db.DB.Exec("INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)", task.Title, task.Description, task.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func GetTasks(c *gin.Context) {

	title := c.Query("title")

	var rows *sql.Rows
	var err error

	if title != "" {
		query := "SELECT id, title, description, status, created_at, completed_at FROM tasks WHERE title LIKE ?"
		likeTitle := "%" + title + "%"
		log.Println("Ejecutando query con filtro:", likeTitle)
		rows, err = db.DB.Query(query, likeTitle)
	} else {
		query := "SELECT id, title, description, status, created_at, completed_at FROM tasks"
		log.Println("Ejecutando query sin filtro")
		rows, err = db.DB.Query(query)
	}

	if err != nil {
		log.Println("Error en la query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch tasks",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			log.Println("Error escaneando fila:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error leyendo datos",
				"details": err.Error(),
			})
			return
		}
		tasks = append(tasks, task)
	}

	log.Println("Tareas encontradas:", len(tasks))
	c.JSON(http.StatusOK, tasks)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si la tarea existe en la base de datos
	var existingTask models.Task
	err := db.DB.QueryRow("SELECT id FROM tasks WHERE id = ?", id).Scan(&existingTask.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Si no se encuentra la tarea, devolver un 404 Not Found
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			// Si ocurre otro error con la base de datos, devolver un 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Validar que el estado sea uno de los permitidos
	validStatuses := []string{"new", "ongoing", "completed"}
	isValidStatus := false
	for _, status := range validStatuses {
		if task.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status, allowed values are: new, ongoing, completed"})
		return
	}

	if task.Status == "completed" {
		task.MarkAsCompleted()
	}

	_, err = db.DB.Exec("UPDATE tasks SET title = ?, description = ?, status = ?, completed_at = ? WHERE id = ?", task.Title, task.Description, task.Status, task.CompletedAt, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	_, err := db.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tarea eliminada correctamente",
		"id":      id,
	})
}
