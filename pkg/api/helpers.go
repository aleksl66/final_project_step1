package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON сериализует data в JSON и пишет в ответ.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "ошибка сериализации ответа", http.StatusInternalServerError)
	}
}
