package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка получения задач: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
