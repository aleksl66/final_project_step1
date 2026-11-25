package api

import (
	"encoding/json"
	"net/http"

	"go_final_project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req addTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON: "+err.Error())
		return
	}

	if req.ID == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "не указан заголовок задачи")
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
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(task); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
