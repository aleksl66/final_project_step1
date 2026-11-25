package api

import (
	"encoding/json"
	"net/http"
	"os"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON: "+err.Error())
		return
	}

	expected := os.Getenv("TODO_PASSWORD")
	if expected == "" {
		writeError(w, http.StatusInternalServerError, "аутентификация не настроена")
		return
	}
	if req.Password != expected {
		writeError(w, http.StatusUnauthorized, "неверный пароль")
		return
	}

	token, err := makeToken(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка создания токена: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
