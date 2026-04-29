package handlers

import (
	"VGarage/internal/database"
	"VGarage/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка БД: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Машина добавлена"})
}
