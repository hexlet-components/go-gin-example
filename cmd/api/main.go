package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hexlet-components/go-gin-example/handlers"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Port   string
	DBPath string
}

func main() {
	var cfg Config

	// Парсинг флагов командной строки
	flag.StringVar(&cfg.Port, "port", "8080", "Port to run the server on")
	flag.StringVar(&cfg.DBPath, "db", "app.db", "Path to SQLite database file")
	flag.Parse()

	// Проверяем существование базы данных
	if _, err := os.Stat(cfg.DBPath); os.IsNotExist(err) {
		log.Fatalf("Database file does not exist: %s\nPlease run migrations first: go run cmd/migrate/main.go up", cfg.DBPath)
	}

	// Подключение к базе данных
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения с БД
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Настройка роутера
	r := handlers.SetupRouter(db)

	// Запуск сервера
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
