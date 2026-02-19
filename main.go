package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vedanthanekar45/novlnest-server/api"
	"github.com/vedanthanekar45/novlnest-server/db"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	db.ConnectDB()
	defer db.Conn.Close()

	queries := db.NewQueries(db.Conn)

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8000/api/v1/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	apiCfg := &api.ApiConfig{
		Store: queries,
		OAuth: oauthConfig,
	}

	router := http.NewServeMux()

	// the good ol' testing routes
	router.HandleFunc("GET /api/v1/health", handleHealth)
	router.HandleFunc("GET /api/v1/users/me", apiCfg.MiddlewareAuth(apiCfg.GetUserData))

	// the fabled authentication routes..
	router.HandleFunc("GET /api/v1/auth/google/login", apiCfg.HandleGoogleLogin)
	router.HandleFunc("GET /api/v1/auth/google/callback", apiCfg.HandleGoogleCallback)

	// the prophecised book routes

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

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
