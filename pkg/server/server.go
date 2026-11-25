package server

import (
	"fmt"
	"net/http"

	"go_final_project/pkg/api"
)

func Start(port string) error {
	webDir := "./web"

	mux := http.NewServeMux()
	api.Init(mux)
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":" + port
	fmt.Printf("Сервер запущен на http://localhost:%s\n", port)
	return http.ListenAndServe(addr, mux)
}
