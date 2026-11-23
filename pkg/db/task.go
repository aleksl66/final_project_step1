package db

import (
	"database/sql"
	"fmt"
	"time"
)

// DeleteTask удаляет задачу по идентификатору.
func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// UpdateDate обновляет только дату задачи.
func UpdateDate(next string, id string) error {
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// GetTask возвращает задачу по её идентификатору.
func GetTask(id string) (*Task, error) {
	t := &Task{}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, err
	}
	return t, nil
}

// UpdateTask обновляет задачу в таблице scheduler.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// Task описывает задачу планировщика.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает её идентификатор.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Tasks возвращает ближайшие задачи, отсортированные по дате.
// search — необязательная строка поиска (по title/comment или по дате dd.mm.yyyy).
func Tasks(limit int, search string) ([]*Task, error) {
	var (
		query string
		args  []any
	)

	if search == "" {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		args = []any{limit}
	} else if d, ok := parseSearchDate(search); ok {
		// Поиск по конкретной дате
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		args = []any{d, limit}
	} else {
		// Поиск по подстроке в title или comment
		like := "%" + search + "%"
		query = `SELECT id, date, title, comment, repeat FROM scheduler
		         WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
		args = []any{like, like, limit}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0) // не nil, чтобы в JSON было []
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// parseSearchDate проверяет, является ли строка датой в формате 02.01.2006,
// и возвращает её в формате 20060102.
func parseSearchDate(s string) (string, bool) {
	t, err := time.Parse("02.01.2006", s)
	if err != nil {
		return "", false
	}
	return t.Format("20060102"), true
}
