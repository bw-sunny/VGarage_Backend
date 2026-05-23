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

	// 1. Достаем РЕАЛЬНЫЕ данные автомобиля из базы данных
	var car models.Car
	carQuery := `SELECT brand, model, year, mileage FROM cars WHERE id = $1`
	err := database.DB.Get(&car, carQuery, req.CarID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Автомобиль не найден в базе данных"})
		return
	}

	// 2. Достаем РЕАЛЬНУЮ историю обслуживания (последние 5 записей, чтобы не раздувать промпт)
	var logs []models.MaintenanceLog
	logsQuery := `SELECT work_date, description, mileage FROM maintenance_logs WHERE car_id = $1 ORDER BY work_date DESC LIMIT 5`
	err = database.DB.Select(&logs, logsQuery, req.CarID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении истории ТО: " + err.Error()})
		return
	}

	fmt.Println(logs)

	// 3. Контекст-буфер (Augmentation) — собираем реальные данные для ИИ
	contextBuffer := fmt.Sprintf("КОНТЕКСТ АВТОМОБИЛЯ:\nМарка/Модель: %s %s (%d г.в.)\nТекущий пробег: %d км\n\n", car.Brand, car.Model, car.Year, car.Mileage)
	contextBuffer += "ИСТОРИЯ ОБСЛУЖИВАНИЯ (Последние записи):\n"

	if len(logs) == 0 {
		contextBuffer += "- Записи об обслуживании отсутствуют.\n"
	} else {
		for _, log := range logs {
			// Красиво форматируем дату без лишнего времени
			formattedDate := log.WorkDate.Format("2006-01-02")
			contextBuffer += fmt.Sprintf("- [%s] %s\n", formattedDate, log.Description)
		}
	}
	contextBuffer += fmt.Sprintf("\nЗАПРОС ПОЛЬЗОВАТЕЛЯ И СИМПТОМЫ:\n%s", req.Message)

	// Печатаем собранный промпт в консоль для отладки
	fmt.Println("================ [ПОДГОТОВЛЕННЫЙ ПРОМПТ ИЗ РЕАЛЬНОЙ БД] ================")
	fmt.Println(contextBuffer)
	fmt.Println("========================================================================")

	// 4. Заглушка ответа (берет данные из настоящей БД)
	lastServiceInfo := "Записей нет"
	if len(logs) > 0 {
		lastServiceInfo = fmt.Sprintf("%s (%s)", logs[0].Description, logs[0].WorkDate.Format("2006-01-02"))
	}

	mockAIResponse := gin.H{
		"status":  "success",
		"is_mock": true,
		"message": fmt.Sprintf(
			"Привет! Я ИИ-ассистент vGarage. Я изучил параметры твоего %s %s (%d г.). "+
				"Последнее зафиксированное действие в сервисной книжке: %s. "+
				"По твоему запросу («%s») я подготавливаю анализ. База полностью настроена, "+
				"промпт сформирован на реальных данных и готов к отправке в DeepSeek!",
			car.Brand, car.Model, car.Year, lastServiceInfo, req.Message,
		),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, mockAIResponse)
}
