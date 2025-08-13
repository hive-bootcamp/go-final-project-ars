package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Arsadidas/go-final-project/pkg/api"
	"github.com/Arsadidas/go-final-project/tests"
	"github.com/go-chi/chi/v5"
)

func Start() {
	webDir := os.Getenv("TODO_WEB_DIR")
	port := tests.Port
	r := chi.NewRouter()
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		parsedPort, err := strconv.Atoi(envPort)
		if err != nil {
			fmt.Printf("Invalid port number in TODO_PORT: %v\n", err)
			return
		}
		port = parsedPort
	}

	fmt.Printf("Server listening on :%d\n", port)

	currentPort := fmt.Sprintf(":%d", port)
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	api.Init(r)
	err := http.ListenAndServe(currentPort, r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
