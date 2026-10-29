package main

import (
	"log"
	"net/http"
	"time"

	"github.com/App-Akademie-Students/ols-fs-010-spec-driven-go/todo-app/backend/internal/httpapi"
	"github.com/App-Akademie-Students/ols-fs-010-spec-driven-go/todo-app/backend/internal/todo"
)

func main() {
	service := todo.NewService(todo.NewInMemoryRepository())
	mux := http.NewServeMux()
	mux.Handle("/api/", httpapi.NewHandler(service))
	mux.Handle("/", http.FileServer(http.Dir("../frontend")))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Todo app listening on http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
