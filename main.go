package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/health", handleHealth)
	router.HandleFunc("POST /api/v1/users", handleCreateUser)
	router.HandleFunc("GET /api/v1/users/{id}", handleGetUser)

	stack := LoggerMiddleware(router)

	server := &http.Server{
		Addr:         ":8000",
		Handler:      stack,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server starting on 8000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("Fetching user " + id))
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User created"))
}
