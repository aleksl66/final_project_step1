package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

// TasksResp — ответ GET /api/tasks.
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler возвращает список ближайших задач.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка получения задач: " + err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
