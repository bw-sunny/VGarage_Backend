package handlers

import (
	"VGarage/internal/database"
	"VGarage/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO users (email, phone, nickname, city)
              VALUES (:email, :phone, :nickname, :city)
			  RETURNING id, language, current_theme, created_at`

	_, err := database.DB.NamedExec(query, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Новый пользователь!"})
}

func CreateCarHandler(c *gin.Context) {
	var car models.Car

	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO cars (
		user_id, is_main, brand, model, year, color, 
		body_type, engine_type, engine_volume, transmission, mileage
	) VALUES (
		:user_id, :is_main, :brand, :model, :year, :color, 
		:body_type, :engine_type, :engine_volume, :transmission, :mileage
	)`

	_, err := database.DB.NamedExec(query, car)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Машина добавлена"})
}

func CreateExpenseHandler(c *gin.Context) {
	var expense models.Expense

	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO expenses (car_id, type, amount)
	VALUES (:car_id, :type, :amount )`

	_, err := database.DB.NamedExec(query, expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Запись в расходы добавлена"})
}
