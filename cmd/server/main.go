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

	v1 := r.Group("/vgarage")
	{
		v1.POST("/create/user", handlers.CreateUserHandler)
		v1.POST("/create/car", handlers.CreateCarHandler)
		v1.POST("/create/expense", handlers.CreateExpenseHandler)

		v1.GET("/user/:id/cars", handlers.GetUserCarsHandler)

		v1.DELETE("/car/:id", handlers.DeleteCarHandler)
	}

	log.Println("Сервер запущен на http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
