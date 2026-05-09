package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           int       `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	Phone        string    `db:"phone" json:"phone"`
	Nickname     string    `db:"nickname" json:"nickname"`
	City         string    `db:"city" json:"city"`
	Language     string    `db:"language" json:"language"`
	CurrentTheme string    `db:"current_theme" json:"current_theme"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Car struct {
	ID           int       `db:"id" json:"id"`
	UserID       int       `db:"user_id" json:"user_id" binding:"required"`
	IsMain       bool      `db:"is_main" json:"is_main"`
	Brand        string    `db:"brand" json:"brand" binding:"required"`
	Model        string    `db:"model" json:"model" binding:"required"`
	Year         int       `db:"year" json:"year"`
	Color        string    `db:"color" json:"color"`
	BodyType     string    `db:"body_type" json:"body_type"`
	EngineType   string    `db:"engine_type" json:"engine_type"`
	EngineVolume float64   `db:"engine_volume" json:"engine_volume"`
	Transmission string    `db:"transmission" json:"transmission"`
	Mileage      int       `db:"mileage" json:"mileage"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Expense struct {
	ID        int       `db:"id" json:"id"`
	CarID     int       `db:"car_id" json:"car_id" binding:"required"`
	Type      string    `db:"type" json:"type" binding:"required"`
	Amount    float64   `db:"amount" json:"amount" binding:"required"`
	Date      time.Time `db:"date" json:"date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Diagnostic struct {
	ID                   int             `db:"id" json:"id"`
	CarID                int             `db:"car_id" json:"car_id" binding:"required"`
	Date                 time.Time       `db:"date" json:"date"`
	ErrorCodes           string          `db:"error_codes" json:"error_codes"`
	Severity             string          `db:"severity" json:"severity"`
	AIResponseStructured json.RawMessage `db:"ai_response_structured" json:"ai_response_structured"`
	CreatedAt            time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time       `db:"updated_at" json:"updated_at"`
}

type AIHistory struct {
	ID                 int             `db:"id" json:"id"`
	UserID             int             `db:"user_id" json:"user_id" binding:"required"`
	CarID              *int            `db:"car_id" json:"car_id"` // *int позволяет записывать nil (null в БД)
	Message            string          `db:"message" json:"message" binding:"required"`
	Response           string          `db:"response" json:"response"`
	ResponseStructured json.RawMessage `db:"response_structured" json:"response_structured"`
	CreatedAt          time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at" json:"updated_at"`
}
