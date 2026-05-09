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

func GetUserCarsHandler(c *gin.Context) {
	userID := c.Param("id")

	var cars []models.Car

	query := `SELECT * FROM cars WHERE user_id = $1 ORDER BY created_at DESC`

	err := database.DB.Select(&cars, query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка машин: " + err.Error()})
		return
	}

	if cars == nil {
		cars = []models.Car{}
	}

	c.JSON(http.StatusOK, cars)
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

func DeleteCarHandler(c *gin.Context) {
	id := c.Param("id")

	// Удаляем машину. Если в базе настроено ON DELETE CASCADE,
	// связанные расходы удалятся сами.
	query := `DELETE FROM cars WHERE id = $1`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить: " + err.Error()})
		return
	}

	// Проверяем, была ли вообще такая машина
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Машина с таким ID не найдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Машина удалена"})
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
