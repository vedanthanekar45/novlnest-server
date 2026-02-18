package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/vedanthanekar45/novlnest-server/api"
	"github.com/vedanthanekar45/novlnest-server/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	db.ConnectDB()
	defer db.Conn.Close()

	queries := db.NewQueries(db.Conn)

	apiCfg := &api.ApiConfig{
		Store: queries,
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/health", handleHealth)

	router.HandleFunc("GET /api/v1/auth/google/login", apiCfg.HandleGoogleLogin)
	router.HandleFunc("GET /api/v1/auth/google/callback", apiCfg.HandleGetUser)

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
