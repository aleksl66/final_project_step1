package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/signin", method(http.MethodPost, signinHandler))
	mux.HandleFunc("/api/nextdate", method(http.MethodGet, nextDateHandler))
	mux.HandleFunc("/api/task", auth(taskHandler))
	mux.HandleFunc("/api/task/done", auth(method(http.MethodPost, doneTaskHandler)))
	mux.HandleFunc("/api/tasks", auth(method(http.MethodGet, tasksHandler)))
}

func method(allowed string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowed {
			http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
