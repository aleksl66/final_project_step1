package api

import (
	"encoding/json"
	"net/http"

	"go_final_project/pkg/db"
)

// updateTaskHandler обновляет задачу.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req addTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON: " + err.Error()})
		return
	}

	if req.ID == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}
	if req.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	task := &db.Task{
		ID:      req.ID,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := checkDate(task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}
