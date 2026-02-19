package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/vedanthanekar45/novlnest-server/db/database"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Getting the nitty-gritty details and credentials..
var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SERVICE"),

	// this is the url the user will be redirected to after they complete or deny the login
	RedirectURL: "http://localhost:8000/api/v1/auth/google/callback",

	// the scope is basically what will we be fetching from the user, in case we will need only email and profile
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

// this is where we handle the login..
func (cfg *ApiConfig) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauth_state := generateStateOauthCookie(w)
	url := cfg.OAuth.AuthCodeURL(oauth_state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (cfg *ApiConfig) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauth_state, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "Missing oauth state cookie", http.StatusUnauthorized)
		return
	}

	if r.FormValue("state") != oauth_state.Value {
		http.Error(w, "Invalid google oauth state", http.StatusUnauthorized)
		return
	}

	code := r.FormValue("code")
	token, err := cfg.OAuth.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "User data fetch failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		http.Error(w, "JSON parsing failed", http.StatusInternalServerError)
		return
	}

	user, err := cfg.Store.CreateUserOrUpdate(r.Context(), db.CreateUserOrUpdateParams{
		Email:     googleUser.Email,
		GoogleID:  pgtype.Text{String: googleUser.ID, Valid: true},
		Name:      pgtype.Text{String: googleUser.Name, Valid: true},
		AvatarUrl: pgtype.Text{String: googleUser.Picture, Valid: true},
	})
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"iss": "novlnest",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := jwtToken.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Signing Error", http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, 200, struct {
		Token string      `json:"token"`
		User  interface{} `json:"user"`
	}{
		Token: tokenString,
		User:  user,
	})

	fmt.Fprintf(w, "Welcome %s! You are logged in.", user.Name.String)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauth_state", Value: state, HttpOnly: true, MaxAge: 360}
	http.SetCookie(w, &cookie)
	return state
}
