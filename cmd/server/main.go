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

	v1 := r.Group("/vgarage/create")
	{
		v1.POST("/user", handlers.CreateUserHandler)
		v1.POST("/car", handlers.CreateCarHandler)
		v1.POST("/expense", handlers.CreateExpenseHandler)
	}

	log.Println("Сервер запущен на http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
