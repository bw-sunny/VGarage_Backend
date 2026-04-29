package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	dsn := "host=localhost port=5432 user=user password=password dbname=vGarage sslmode=disable"

	var err error

	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	log.Println("--- Подключение к PostgreSQL успешно установлено! ---")
}
