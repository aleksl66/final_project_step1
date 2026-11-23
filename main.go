package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	// Инициализация базы данных
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close()

	// Определяем порт
	port := "7540"
	if p := os.Getenv("TODO_PORT"); p != "" {
		port = p
	}

	// Запуск сервера
	if err := server.Start(port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
