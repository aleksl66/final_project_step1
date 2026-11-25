package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// schema — SQL-команды для создания таблицы и индекса.
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8) NOT NULL DEFAULT "",
    title   VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat  VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// db — глобальное подключение к базе данных.
var db *sql.DB

// Init открывает базу данных dbFile и при необходимости создаёт таблицу.
func Init(dbFile string) error {
	// Проверяем, существует ли файл базы
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем базу данных
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Если файла не было — создаём таблицу и индекс
	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы: %w", err)
		}
	}

	return nil
}

// GetDB возвращает глобальное подключение к базе данных.
func GetDB() *sql.DB {
	return db
}

// Close закрывает подключение к базе данных.
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
