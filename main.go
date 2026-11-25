package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	// Путь к БД: TODO_DBFILE или по умолчанию
	dbFile := "scheduler.db"
	if f := os.Getenv("TODO_DBFILE"); f != "" {
		dbFile = f
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	port := "7540"
	if p := os.Getenv("TODO_PORT"); p != "" {
		port = p
	}

	if err := server.Start(port); err != nil {
		db.Close()
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
	db.Close()
}
