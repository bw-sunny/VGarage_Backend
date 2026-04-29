package main

import (
	"VGarage/internal/database"
	"VGarage/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()

	// Группа API v1
	v1 := r.Group("/api/v1")
	{
		v1.POST("/cars", handlers.CreateCarHandler)
		// Сюда добавим POST /users, POST /expenses и т.д.
	}

	log.Println("Сервер стартовал на :8080")
	r.Run(":8080")
}
