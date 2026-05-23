package handlers

import (
	"VGarage/internal/database"
	"VGarage/internal/models"
	"VGarage/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// ==========================================
// СТРУКТУРЫ ДАННЫХ (Для Карт и ИИ)
// ==========================================

// Структура ответа для фронтенда (Поиск на карте)
type SearchResult struct {
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Coordinates []float64 `json:"coordinates"` // [lon, lat]
}

// Облегченная структура ответа Yandex Geosearch API
type YandexGeosearchResponse struct {
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"properties"`
	} `json:"features"`
}

// RequestAI определяет то, что прилетает с React-фронтенда для ИИ
type RequestAI struct {
	CarID   int    `json:"car_id" binding:"required"`
	Message string `json:"message" binding:"required"` // Симптомы или вопрос
}

// Модели для имитации базы данных ИИ
type CarSpec struct {
	Brand string
	Model string
	Year  int
}

type MaintenanceLog struct {
	Date        string
	Description string
}

// ==========================================
// ХЭНДЛЕРЫ ПОЛЬЗОВАТЕЛЕЙ И МАШИН
// ==========================================

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

	query := `DELETE FROM cars WHERE id = $1`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить: " + err.Error()})
		return
	}

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

// ==========================================
// ХЭНДЛЕРЫ ГЕОЛОКАЦИИ И КАРТ
// ==========================================

func MapHandler(c *gin.Context) {
	c.File("web/map.html")
}

func GeocodeHandler(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр address обязателен"})
		return
	}

	coordinates, err := services.FetchCoordinates(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address":     address,
		"coordinates": coordinates,
	})
}

// Поиск организаций с учетом геолокации пользователя
func GeosearchHandler(c *gin.Context) {
	searchText := c.Query("text")
	if searchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр text обязателен"})
		return
	}

	lon := c.Query("lon")
	lat := c.Query("lat")

	geosearchKey := "ТВОЙ_КЛЮЧ_ОТ_API_ГЕОПОИСКА"

	apiURL := fmt.Sprintf(
		"https://search-maps.yandex.ru/v1/?text=%s&key=%s&lang=ru_RU&limit=15",
		url.QueryEscape(searchText),
		geosearchKey,
	)

	if lon != "" && lat != "" {
		apiURL = fmt.Sprintf("%s&ll=%s,%s&spn=0.1,0.1&rspn=0", apiURL, lon, lat)
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка запроса к Яндексу"})
		return
	}
	defer resp.Body.Close()

	var yandexData YandexGeosearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&yandexData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка декодирования данных"})
		return
	}

	var clientResults []SearchResult
	for _, feature := range yandexData.Features {
		if len(feature.Geometry.Coordinates) == 2 {
			clientResults = append(clientResults, SearchResult{
				Name:        feature.Properties.Name,
				Address:     feature.Properties.Description,
				Coordinates: feature.Geometry.Coordinates,
			})
		}
	}

	c.JSON(http.StatusOK, clientResults)
}

// ==========================================
// ХЭНДЛЕР ИИ-АССИСТЕНТА (Псевдо-RAG)
// ==========================================

func AIAssistantHandler(c *gin.Context) {
	var req RequestAI
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// 1. Имитируем быстрый SQL-запрос к таблице `cars`
	car := CarSpec{
		Brand: "Toyota",
		Model: "Camry",
		Year:  2018,
	}

	// 2. Имитируем SQL-запрос к `maintenance_logs`
	logs := []MaintenanceLog{
		{Date: "2026-01-15", Description: "Замена моторного масла и фильтров"},
		{Date: "2025-10-10", Description: "Диагностика подвески, замена передних сайлентблоков"},
	}

	// 3. Контекст-буфер (Augmentation)
	contextBuffer := fmt.Sprintf("КОНТЕКСТ АВТОМОБИЛЯ:\nМарка/Модель: %s %s (%d г.в.)\n", car.Brand, car.Model, car.Year)
	contextBuffer += "ИСТОРИЯ ОБСЛУЖИВАНИЯ:\n"
	for _, log := range logs {
		contextBuffer += fmt.Sprintf("- [%s] %s\n", log.Date, log.Description)
	}
	contextBuffer += fmt.Sprintf("\nЗАПРОС ПОЛЬЗОВАТЕЛЯ:\n%s", req.Message)

	// Печатаем собранный промпт в консоль для отладки
	fmt.Println("================ [ПОДГОТОВЛЕННЫЙ ПРОМПТ ДЛЯ ИИ] ================")
	fmt.Println(contextBuffer)
	fmt.Println("================================================================")

	// 4. Заглушка ответа
	mockAIResponse := gin.H{
		"status":  "success",
		"is_mock": true,
		"message": fmt.Sprintf(
			"Привет! Я твой ИИ-ассистент vGarage. Вижу, что ты владеешь %s %s %d года. "+
				"Последний раз ты менял масло %s. На основе твоего запроса («%s»), я пока "+
				"имитирую анализ. Как только мой создатель подключит API DeepSeek, здесь появится "+
				"реальный диагноз поломки!",
			car.Brand, car.Model, car.Year, logs[0].Date, req.Message,
		),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, mockAIResponse)
}
