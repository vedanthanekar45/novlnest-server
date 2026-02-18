package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/vedanthanekar45/novlnest-server/db"
)

type APIServer struct {
	addr  string
	store *db.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	db.ConnectDB()
	defer db.Conn.Close()

	queries := db.NewQueries(db.Conn)

	app := &APIServer{
		addr:  ":8000",
		store: queries,
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/health", app.handleHealth)
	router.HandleFunc("POST /api/v1/users", app.handleCreateUser)
	router.HandleFunc("GET /api/v1/users/{id}", app.handleGetUser)

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

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("Fetching user " + id))
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User created"))
}
