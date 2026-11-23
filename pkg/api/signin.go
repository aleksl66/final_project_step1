package api

import (
	"encoding/json"
	"net/http"
	"os"
)

// signinHandler проверяет пароль и возвращает JWT.
func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON: " + err.Error()})
		return
	}

	expected := os.Getenv("TODO_PASSWORD")
	if expected == "" {
		writeJSON(w, map[string]string{"error": "аутентификация не настроена"})
		return
	}
	if req.Password != expected {
		writeJSON(w, map[string]string{"error": "неверный пароль"})
		return
	}

	token, err := makeToken(req.Password)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка создания токена: " + err.Error()})
		return
	}

	writeJSON(w, map[string]string{"token": token})
}
